package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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

	passwordHash, err := cryptography.HashPassword(user.PasswordHash, p)
	if err != nil {
		fmt.Printf("[!] ERROR: hashing password %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	otpSecret := auth.GenerateTOTPSecret()
	fmt.Printf("[+] OTP secret: %v\n", otpSecret)
	if otpSecret == "" {
		fmt.Printf("[!] ERROR: generating random secret for otp")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	encryptedStoredBytesOtpSecret := cryptography.EncryptDataRSA([]byte(otpSecret))
	if encryptedStoredBytesOtpSecret == nil {
		fmt.Printf("[!] ERROR: encrypting otp secret to store")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}
	encryptedStoredOtpSecret := cryptography.Base64Encode(encryptedStoredBytesOtpSecret)

	createUser, err := clientPostgresDb.User.CreateOne(
		db.User.Username.Set(user.Username),
		db.User.Email.Set(user.Email),
		db.User.MasterPasswordHash.Set(passwordHash),
		db.User.Salt.Set(user.Salt),
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
	fmt.Printf("password lenght: %v\n", len(passwordHash))
	fmt.Print("password hash: ", passwordHash, "\n")
	// encryptedBytesOtpSecret, err := cryptography.EncryptAESGCM(otpSecretUri, []byte(passwordHash))
	// if err != nil {
	// 	fmt.Printf("[!] ERROR: encrypting otp secret to send back: %v", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"message": "could not create user",
	// 	})
	// }
	// encryptedOtpSecret := cryptography.Base64Encode(encryptedBytesOtpSecret)

	// a := cryptography.Base64Decode(encryptedOtpSecret)
	// z, _ := cryptography.DecryptAESCBC(a, []byte(passwordHash)) 	

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User created successfully",
		// "otpSecretUri": encryptedOtpSecret,
		"otpSecretUriUnencrypted": otpSecretUri,
		// "otpSecretUriDecrypted": string(z),
	})
}

func Salt(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"salt": cryptography.GenerateSalt(p.SaltLength),
	})
}

func Login(c *fiber.Ctx) error {
	var user models.UserLogin

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

	retrivedUserDb, err := clientPostgresDb.User.FindMany(
		db.User.Email.Equals(user.Email),
	).Exec(context.Background())
	if errors.Is(err, db.ErrNotFound) {
		fmt.Printf("[-] No record with email: %s\n", user.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email incorrect",
		})
	} else if err != nil {
		fmt.Printf("[!] Error occurred finding email in database: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// check if the password is correct
	equal, err := cryptography.ComparePasswordAndHash(user.PasswordHash, retrivedUserDb[0].MasterPasswordHash)
	if err != nil {
		fmt.Printf("[!] Error occurred comparing password and hash: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	
	if !equal {
		fmt.Printf("[-] Login: Password wrong\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password incorrect",
		})
	}

	// Check if the user has an OTP secret
	otpSecret := clientRedisDb.Get(context.Background(), retrivedUserDb[0].ID).Val()
	if otpSecret == "" {
		fmt.Printf("[-] OTP: No secret found\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "OTP error",
		})
	}
	if !auth.VerifyOTP(otpSecret, user.Otp) {
		fmt.Printf("[-] OTP: Invalid code\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	fmt.Printf("[+] OTP: Access granted\n")

	jwtToken, err := auth.GenerateJWTToken(retrivedUserDb[0].ID, os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		fmt.Printf("[!] Error occurred generating JWT token: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	cookie := fiber.Cookie{
		Name:     os.Getenv("JWT_COOKIE_TOKEN_NAME"),
		Value:    jwtToken,
		Expires:  utils.CalculateExpireTime(os.Getenv("JWT_EXPIRES_IN")),
		HTTPOnly: true,
	}

	_, err = clientPostgresDb.User.FindMany(
		db.User.ID.Equals(retrivedUserDb[0].ID),
	).Update(
		db.User.LastLogin.Set(time.Now()),
	).Exec(context.Background())

	if err != nil {
		fmt.Printf("[!] Error occurred updating last login: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Printf("[+] Updated last login: \n",)

	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

func TestJWT(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "jwt",
	})
}

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
