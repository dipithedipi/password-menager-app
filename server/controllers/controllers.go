package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dipithedipi/password-manager/auth"
	"github.com/dipithedipi/password-manager/cryptography"
	"github.com/dipithedipi/password-manager/models"
	"github.com/dipithedipi/password-manager/prisma/db"
	"github.com/dipithedipi/password-manager/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var clientPostgresDb *db.PrismaClient
var clientRedisDb *redis.Client 
var p = &models.ArgonParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}


func SetPostgresDbClient(client *db.PrismaClient) {
	clientPostgresDb = client
}

func SetRedisDbClient(client *redis.Client) {
	clientRedisDb = client
}

func Register(c *fiber.Ctx) error {
	var user models.UserRegister

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(user) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing required fields",
		})
	}

	if !utils.ValidateEmail(user.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid email",
		})
	}

	ctx := context.Background()

	passwordHash, _, err := cryptography.HashPassword(user.PasswordHash, p)
	if err != nil {
		fmt.Printf("[!] ERROR: hashing password %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	otpSecret := auth.GenerateTOTPSecret()
	fmt.Printf("[+] OTP secret: %v\n", otpSecret)
	if otpSecret == "" {
		fmt.Printf("[!] ERROR: generating random secret for otp: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	encryptedStoredOtpSecret := cryptography.EncryptDataRSA([]byte(otpSecret))
	if encryptedStoredOtpSecret == nil {
		fmt.Printf("[!] ERROR: encrypting otp secret to store: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	createUser, err := clientPostgresDb.User.CreateOne(
		db.User.Username.Set(user.Username),
		db.User.Email.Set(user.Email),
		db.User.MasterPasswordHash.Set(passwordHash),
		db.User.Salt.Set([]byte(user.Salt)),
		db.User.OtpSecret.Set(encryptedStoredOtpSecret),
	).Exec(ctx)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "failed on the fields: (`username`)"):
			fmt.Printf("[!] ERROR REGISTER: Username already exists %v \n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Username already exists",
			})
		case strings.Contains(err.Error(), "failed on the fields: (`email`)"):
			fmt.Printf("[!] ERROR REGISTER: Email already exists %v \n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email already exists",
			})
		default:
			// General database error
			fmt.Printf("[!] ERROR REGISTER: Internal server error %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not create user",
			})
		}
	}

	result, err := json.MarshalIndent(createUser, "", " ")
	if err != nil {
		fmt.Printf("[!] ERROR: Marshalling user data %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not marshal user data",
		})
	}

	clientRedisDb.Set(ctx, createUser.ID, otpSecret, 0)
	fmt.Printf("[+] OTP secret stored in redis(%v:%v)\n", createUser.ID, otpSecret)

	fmt.Printf("[+] Created user: %s\n", result)

	otpSecretUri := auth.GenerateUriTOTPWithSecret(otpSecret, user.Email)
	encryptedBytesOtpSecret, err := cryptography.EncryptAESCBC(otpSecretUri, []byte(passwordHash))
	if err != nil {
		fmt.Printf("[!] ERROR: encrypting otp secret to send back: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}
	encryptedOtpSecret := cryptography.Base64Encode(encryptedBytesOtpSecret)

	// a := cryptography.Base64Decode(encryptedOtpSecret)
	// z, _ := cryptography.DecryptAESCBC(a, []byte(passwordHash)) 	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		"otpSecretUri": encryptedOtpSecret,
		//"otpSecretUri": otpSecretUri,
		//"otpSecretUriDecrypted": string(z),
	})
}

func Salt(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"salt": cryptography.GenerateSalt(p.SaltLength),
	})
}

// const SecretKey = "secret"

// func Login(c *fiber.Ctx) error {
// 	var data map[string]string

// 	if err := c.BodyParser(&data); err != nil {
// 		return err
// 	}

// 	var user models.User

// 	database.DB.Where("email = ?", data["email"]).First(&user) //Check the email is present in the DB

// 	if user.ID == 0 { //If the ID return is '0' then there is no such email present in the DB
// 		c.Status(fiber.StatusNotFound)
// 		return c.JSON(fiber.Map{
// 			"message": "user not found",
// 		})
// 	}

// 	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
// 		c.Status(fiber.StatusBadRequest)
// 		return c.JSON(fiber.Map{
// 			"message": "incorrect password",
// 		})
// 	} // If the email is present in the DB then compare the Passwords and if incorrect password then return error.

// 	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		Issuer:    strconv.Itoa(int(user.ID)),
// 		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
// 	})

// 	token, err := claims.SignedString([]byte(SecretKey))

// 	if err != nil {
// 		c.Status(fiber.StatusInternalServerError)
// 		return c.JSON(fiber.Map{
// 			"message": "could not login",
// 		})
// 	}

// 	cookie := fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    token,
// 		Expires:  time.Now().Add(time.Hour * 24),
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&cookie)

// 	return c.JSON(fiber.Map{
// 		"message": "success",
// 	})
// }

// func User(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")

// 	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(SecretKey), nil
// 	})

// 	if err != nil {
// 		c.Status(fiber.StatusUnauthorized)
// 		return c.JSON(fiber.Map{
// 			"message": "unauthenticated",
// 		})
// 	}

// 	claims := token.Claims.(*jwt.StandardClaims)

// 	var user models.User

// 	database.DB.Where("id = ?", claims.Issuer).First(&user)

// 	return c.JSON(user)

// }

// func Logout(c *fiber.Ctx) error {
// 	cookie := fiber.Cookie{
// 		Name:     "jwt",
// 		Value:    "",
// 		Expires:  time.Now().Add(-time.Hour), //Sets the expiry time an hour ago in the past.
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&cookie)

// 	return c.JSON(fiber.Map{
// 		"message": "success",
// 	})

// }
