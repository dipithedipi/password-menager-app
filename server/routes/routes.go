package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/dipithedipi/password-manager/controllers" // importing the routes package
	"github.com/dipithedipi/password-manager/cryptography"
	"github.com/dipithedipi/password-manager/prisma/db"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func SetupClientPostgresDb () (*db.PrismaClient, error) {
    client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
	  return nil, err
	}

    return client, nil
}

func SetupClientRedisDb () (*redis.Client, error){
    // Connect to the redis database
    db, _ := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"), // Replace with your Redis server address
        Password: os.Getenv("REDIS_PASSWORD"),                // No password for local development
        DB:       db,                 // Default DB
    })
      
    // Ping the Redis server to check the connection
    var ctx = context.Background()
    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, err
    } else {
        return client, nil
    }
}

func SyncRedisOTPSecret(clientPostgresDb *db.PrismaClient, clientRedisDb *redis.Client) {

    // Flush the redis database
    // err := clientRedisDb.FlushDB(context.Background()).Err()
	// if err != nil {
	// 	panic(err)
	// }


    // get all from redis
    keys, err := clientRedisDb.Keys(context.Background(), "*").Result()
    if err != nil {
        panic(err)
    }

    // Get all the users from the postgres database
    users, err := clientPostgresDb.User.FindMany().Exec(context.Background())
    if err != nil {
        panic(err)
    }

    // Check if the redis database is already in sync with the Postgres database
    if len(keys) == len(users) {
        fmt.Println("[+] Redis database is already in sync with the Postgres database")
        return
    }

    // Iterate over the users and sync the OTP secrets with the redis database
    for _, user := range users {
        // Check if the user has an OTP secret
        if user.OtpSecret != "" {
            // Add the OTP secret to the redis database
            encryptedOtpSecretBytes := cryptography.Base64Decode(user.OtpSecret)
            unencryptedOtpSecret := cryptography.DecryptDataRSA(encryptedOtpSecretBytes)

            err = clientRedisDb.Set(context.Background(), user.ID, string(unencryptedOtpSecret), 0).Err()
            if err != nil {
                panic(err)
            }
        }
    }
    fmt.Println("[+] Redis database is now in sync with the Postgres database")
}

func Setup(app *fiber.App, clientPostgresDb *db.PrismaClient, clientRedisDb *redis.Client) {
    controllers.SetPostgresDbClient(clientPostgresDb)
    controllers.SetRedisDbClient(clientRedisDb)

    apiUser := app.Group("/user")
    apiUser.Post("/register", controllers.Register)
    apiUser.Get("/register/salt", controllers.Salt)
    apiUser.Post("/login", controllers.Login)

    app.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
        TokenLookup: "cookie:" + os.Getenv("JWT_COOKIE_TOKEN_NAME"),
    }))

    apiPassword := app.Group("/password")
    apiPassword.Get("/jwt", controllers.TestJWT)

    // api.Get("/get-user", controllers.User)
    // api.Post("/login", controllers.Login)
    // api.Post("/logout", controllers.Logout)
}