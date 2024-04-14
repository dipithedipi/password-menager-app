package main

import (
	"fmt"
	"github.com/dipithedipi/password-manager/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    clientDb, err := routes.SetupClientDb()
    if err != nil {
        fmt.Println("[!] ERROR: Unable to connect to the database")
        panic(err)
    } else {
        fmt.Println("[+] Connected to the database")
    }

    routes.Setup(app, clientDb)

    app.Listen(":8000")
}