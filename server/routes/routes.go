package routes

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/dipithedipi/password-manager/auth"
	"github.com/dipithedipi/password-manager/controllers" // importing the routes package
	"github.com/dipithedipi/password-manager/cryptography"
	"github.com/dipithedipi/password-manager/prisma/db"
	"github.com/dipithedipi/password-manager/utils"
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
            unencryptedOtpSecret := cryptography.DecryptServerDataRSA(encryptedOtpSecretBytes)

            err = clientRedisDb.Set(context.Background(), user.ID, string(unencryptedOtpSecret), 0).Err()
            if err != nil {
                panic(err)
            }
        }
    }
    fmt.Println("[+] Redis database is now in sync with the Postgres database")
}

func SetupKAnonymity(clientPostgresDb *db.PrismaClient) {
    // check if the database is empty
    result, err := clientPostgresDb.PasswordLeak.FindMany().Take(100).Exec(context.Background())
    if err != nil {
        fmt.Println("[!] Error checking if the password leaked database is empty")
        panic(err)
    }

    if len(result) > 0 {
        fmt.Println("[+] Password leaked database is already populated")
    } else {
        // first time setup
        // check if the file exists or is empty
        _, err := os.Stat("wordlist.tmp");
        if os.IsNotExist(err) {

            // download a file 
            fmt.Println("[+] Downloading the starting wordlist file")
            err := utils.DownloadFile("wordlist.tmp", os.Getenv("K_ANONYMITY_DEFAULT_WORDLIST_LINK"))
            if err != nil {
                fmt.Println("[!] Error downloading the starting wordlist file")
                panic(err)
            }
        }

        // read the file
        lines, err := utils.ReadFileContent("wordlist.tmp")
        if err != nil {
            fmt.Println("[!] Error reading the starting wordlist file")
            panic(err)
        }

        // check if the file is empty
        if len(lines) == 0 {
            fmt.Println("[!] Error: The starting wordlist file is empty")
            panic(err)
        }

        // insert the data into the database
        for _, line := range lines {
            _, err := clientPostgresDb.PasswordLeak.CreateOne(
                db.PasswordLeak.PasswordHash.Set(cryptography.Sha1(line)),
            ).Exec(context.Background())
            if err != nil {
                fmt.Println("[!] Error inserting the starting wordlist file into the database")
                panic(err)
            }
        }

        // delete the file
        err = os.Remove("wordlist.tmp")
        if err != nil {
            fmt.Println("[!] Error deleting the starting wordlist file")
            panic(err)
        }

        fmt.Println("[+] Password leaked database is now populated")
    }
}

func Setup(app *fiber.App, clientPostgresDb *db.PrismaClient, clientRedisDb *redis.Client) {
    controllers.SetPostgresDbClient(clientPostgresDb)
    controllers.SetRedisDbClient(clientRedisDb)

    apiUser := app.Group("/user")
    apiUser.Post("/register", controllers.Register)
    apiUser.Get("/register/salt", controllers.RandomSalt)
    apiUser.Post("/login", controllers.Login)
    apiUser.Get("/login/salt", controllers.GetSaltFromUser)
    apiUser.Post("/logout", controllers.Logout)

    apiPassword := app.Group("/password", auth.MiddlewareJWTAuth(clientRedisDb))
    apiPassword.Post("/new", controllers.PostNewPassword)
    apiPassword.Get("/search", controllers.GetPasswordPreview)
    apiPassword.Get("/get", controllers.GetPassword)

    apiUtils := app.Group("/utils")
    apiUtils.Get("/checkPassword", controllers.CheckPasswordLeak)

    // api.Get("/get-user", controllers.User)
}