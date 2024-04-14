package routes

import (
    "github.com/dipithedipi/password-manager/controllers" // importing the routes package 
    "github.com/gofiber/fiber/v2"
    "github.com/dipithedipi/password-manager/prisma/db"
)

func SetupClientDb () (*db.PrismaClient, error) {
    client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
	  return nil, err
	}

    return client, nil
}

func Setup(app *fiber.App, clientDb *db.PrismaClient) {
    controllers.SetDbClient(clientDb)

    api := app.Group("/user")
    api.Post("/register", controllers.Register)
    api.Get("/register/salt", controllers.Salt)
    // api.Get("/get-user", controllers.User)
    // api.Post("/login", controllers.Login)
    // api.Post("/logout", controllers.Logout)
}