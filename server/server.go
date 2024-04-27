package main

import (
	"fmt"
    "strconv"
	"os"
	"github.com/dipithedipi/password-manager/cryptography"
	"github.com/dipithedipi/password-manager/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
    // Load the .env file
    err := godotenv.Load()
	if err != nil {
		fmt.Println("[!] ERROR: Loading .env file")
		return
	} else {
        fmt.Println("[+] Loaded .env file")
    }

    // Generate encryption keys if they don't exist
    if _, err := os.Stat(os.Getenv("PUBLIC_KEY_PATH")); os.IsNotExist(err) {
        fmt.Println("[!] Generating crypto keys")
        keyLenght, _ := strconv.Atoi(os.Getenv("RSA_KEY_LENGTH"))
        cryptography.GenerateKeysRSA(os.Getenv("PUBLIC_KEY_PATH"), os.Getenv("PRIVATE_KEY_PATH"), keyLenght)
        fmt.Print("[+] Keys generated and saved\n")
    } else {
        fmt.Println("[+] Crypto keys already exist")
    }

    // Connect to the Postgres database
    clientPostgresDb, err := routes.SetupClientPostgresDb()
    if err != nil {
        fmt.Println("[!] ERROR: Unable to connect to the postgres database")
        panic(err)
    } else {
        fmt.Println("[+] Connected to the postgres database")
    }

    // Connect to the redis database
    clientRedisDb, err := routes.SetupClientRedisDb()
    if err != nil {
        fmt.Println("[!] ERROR: Unable to connect to the redis database")
        panic(err)
    } else {
        fmt.Println("[+] Connected to the redis database")
    }

    // Sync the redis database with otp secrets decrypted from the postgres database
    routes.SyncRedisOTPSecret(clientPostgresDb, clientRedisDb)

    // Create a new Fiber app
    app := fiber.New()
    routes.Setup(app, clientPostgresDb, clientRedisDb)
    app.Listen(":8000")
}