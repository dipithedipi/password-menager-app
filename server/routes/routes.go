package routes

import (
    "github.com/dipithedipi/password-manager/controllers" // importing the routes package 
    "github.com/gofiber/fiber/v2"
    "github.com/dipithedipi/password-manager/prisma/db"
    "strconv"
    "github.com/redis/go-redis/v9"
    "os"
    "context"
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


func Setup(app *fiber.App, clientPostgresDb *db.PrismaClient, clientRedisDb *redis.Client) {
    controllers.SetPostgresDbClient(clientPostgresDb)
    controllers.SetRedisDbClient(clientRedisDb)

    api := app.Group("/user")
    api.Post("/register", controllers.Register)
    api.Get("/register/salt", controllers.Salt)
    // api.Get("/get-user", controllers.User)
    // api.Post("/login", controllers.Login)
    // api.Post("/logout", controllers.Logout)
}