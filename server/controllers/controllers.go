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
	"github.com/dipithedipi/password-manager/event"
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
	SaltLength:  32,
	KeyLength:   64,
}

func SetPostgresDbClient(client *db.PrismaClient) {
	clientPostgresDb = client
}

func SetRedisDbClient(client *redis.Client) {
	clientRedisDb = client
}

// ------------------- api controllers -------------------

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

	encryptedStoredBytesOtpSecret, err := cryptography.EncryptServerDataRSA([]byte(otpSecret))
	if err != nil {
		fmt.Printf("[!] ERROR: encrypting otp secret to store")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}
	encryptedStoredOtpSecret := cryptography.Base64Encode(encryptedStoredBytesOtpSecret)

	otpSecretUri := auth.GenerateUriTOTPWithSecret(otpSecret, user.Email)

	// convert public key to pem
	publicKeyPEM, err := cryptography.ConvertBase64PublicKeyToPEM(user.PublicKey)
	if err != nil {
		fmt.Printf("[!] ERROR: converting public key to pem %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "public key not valid",
		})
	}

	encOtpSecretUriBytes, err := cryptography.EncryptDataRSA([]byte(otpSecretUri), []byte(publicKeyPEM))
	if err != nil {
		fmt.Printf("[!] ERROR: encrypting otp secret to send back: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create user",
		})
	}

	encOtpSecreturi := cryptography.Base64Encode(encOtpSecretUriBytes)

	// add to db
	createUser, err := clientPostgresDb.User.CreateOne(
		db.User.Username.Set(user.Username),
		db.User.Email.Set(user.Email),
		db.User.MasterPasswordHash.Set(passwordHash),
		db.User.Salt.Set(user.Salt),
		db.User.OtpSecret.Set(encryptedStoredOtpSecret),
		db.User.PublicKey.Set(user.PublicKey),
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

	// add new event do user db

	err = event.NewEvent(clientPostgresDb, "User created", "", c.IP(), createUser.ID)
	if err != nil {
		fmt.Printf("[!] ERROR: creating event %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "User created successfully",
		"otpSecretUri": encOtpSecreturi,
		// "otpSecretUriUnencrypted": otpSecretUri,
		// "otpSecretUriDecrypted": string(z),
	})
}

func GetSaltFromUser(c *fiber.Ctx) error {
	var user models.UserSaltLogin

	// read from url params
	if err := c.BodyParser(&user); err != nil {
		fmt.Printf("error: %v\n", err)
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
	).Select(
		db.User.Salt.Field(),
	).Exec(context.Background())
	if errors.Is(err, db.ErrNotFound) {
		fmt.Printf("[-] No record for salt with email: %s\n", user.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email incorrect",
		})
	} else if err != nil {
		fmt.Printf("[!] Error occurred finding salt from email in database: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(retrivedUserDb) == 0 {
		fmt.Printf("[-] No record for salt with email: %s\n", user.Email)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Email not found",
		})
	}

	salt := retrivedUserDb[0].Salt
	return c.JSON(fiber.Map{
		"salt": salt,
	})
}

func RandomSalt(c *fiber.Ctx) error {
	rawSalt, err := cryptography.GenerateSalt(p.SaltLength)
	if err != nil {
		fmt.Printf("[!] ERROR: generating salt %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not generate salt",
		})
	}

	salt := cryptography.Base64Encode(rawSalt)

	return c.JSON(fiber.Map{
		"salt": salt,
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
		// Event login failed
		err = event.NewEvent(clientPostgresDb, "Login failed", "User used a wrong password", c.IP(), retrivedUserDb[0].ID)
		if err != nil {
			fmt.Printf("[!] Error occurred creating event: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

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

		// Event login failed
		err = event.NewEvent(clientPostgresDb, "Login failed", "User use an OTP code invalid", c.IP(), retrivedUserDb[0].ID)
		if err != nil {
			fmt.Printf("[!] Error occurred creating event: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	fmt.Printf("[+] OTP: Access granted\n")

	jwtToken, err := auth.GenerateJWTToken(retrivedUserDb[0].ID, c.IP(), os.Getenv("JWT_EXPIRES_IN"))
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

	// delete old token(pass expire date) from postgres db for sessions
	_, err = clientPostgresDb.Token.FindMany(
		db.Token.UserID.Equals(retrivedUserDb[0].ID),
		db.Token.ExpireAt.Before(time.Now()),
	).Delete().Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Session: Error occurred deleting old tokens, %s\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// add token to postgres db for sessions online
	_, err = clientPostgresDb.Token.CreateOne(
		db.Token.TokenValue.Set(jwtToken),
		db.Token.ExpireAt.Set(utils.CalculateExpireTime(os.Getenv("JWT_EXPIRES_IN"))),
		db.Token.IPAddress.Set(c.IP()),
		db.Token.User.Link(
			db.User.ID.Equals(retrivedUserDb[0].ID),
		),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Session: Error occurred adding token to db: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// Event login success
	err = event.NewEvent(clientPostgresDb, "Login success", "User logged in", c.IP(), retrivedUserDb[0].ID)
	if err != nil {
		fmt.Printf("[!] Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	fmt.Printf("[+] Updated last login: \n")
	fmt.Printf("[+] Session: new session\n")
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

func PostNewPassword(c *fiber.Ctx) error {
	// get user id from jwt
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "jwt error",
		})
	}

	var passwordFields models.PasswordSet
	if err := c.BodyParser(&passwordFields); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordFields) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing required fields",
		})
	}

	// add password linked to user id in DB
	_, err = clientPostgresDb.Password.CreateOne(
		db.Password.Website.Set(passwordFields.Domain),
		db.Password.Username.Set(passwordFields.Username),
		db.Password.Password.Set(passwordFields.Password),
		db.Password.User.Link(
			db.User.ID.Equals(claims.Issuer),
		),
		db.Password.Description.Set(passwordFields.Description),
		db.Password.OtpProtected.Set(passwordFields.Otp),
	).Exec(context.Background())

	if err != nil {
		fmt.Printf("[!] Error occurred adding password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error adding password",
		})
	}

	// event password added
	err = event.NewEvent(clientPostgresDb, "Password added", fmt.Sprintf("User added a new password for %s", passwordFields.Domain), c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	fmt.Print("[+] Password added successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password added successfully",
	})
}

func GetPasswordPreview(c *fiber.Ctx) error {
	// get user id from jwt
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "jwt error",
		})
	}

	var passwordRequest models.PasswordRequestSearch
	if err := c.BodyParser(&passwordRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing required fields",
		})
	}

	// get password field without the password linked to user id in DB
	result, err := clientPostgresDb.Password.FindMany(
		db.Password.Website.Contains(passwordRequest.Domain),
		db.Password.UserID.Equals(claims.Issuer),
	).Omit(
		db.Password.Password.Field(),
		db.Password.UserID.Field(),
	).Exec(context.Background())

	if err != nil {
		fmt.Printf("[!] Error occurred finding password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding password",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] No password found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No password found",
		})
	}

	// Convertiamo []db.PasswordModel in []interface{}
	var interfaceSlice []interface{}
	for _, password := range result {
		interfaceSlice = append(interfaceSlice, password)
	}

	fieldsToRemove := []string{"userId", "password"}
	clearedResult, err := utils.ClearJsonFields(interfaceSlice, fieldsToRemove)
	if err != nil {
		fmt.Printf("[!] Error occurred clearing json fields: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error elaborating password",
		})
	}

	// event password preview
	err = event.NewEvent(clientPostgresDb, "Password preview", fmt.Sprintf("User previewed a passwords for %s", passwordRequest.Domain), c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	fmt.Printf("[+] Password found:")
	utils.PrintFormattedJSON(clearedResult)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Password found",
		"passwords": clearedResult,
	})
}

func GetPassword(c *fiber.Ctx) error {
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	var passwordRequest models.PasswordRequestInfo
	if err := c.BodyParser(&passwordRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
		})
	}

	result, err := clientPostgresDb.Password.FindMany(
		db.Password.ID.Equals(passwordRequest.PasswordId),
		db.Password.Website.Contains(passwordRequest.Domain),
		db.Password.UserID.Equals(claims.Issuer),
	).Exec(context.Background())

	if err != nil {
		fmt.Printf("[!] Error occurred finding password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding password",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] No password found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No password found",
		})
	}

	// check if otp is needed
	if result[0].OtpProtected {
		fmt.Printf("[+] Password is protected by OTP\n")
		otpSecret := clientRedisDb.Get(context.Background(), claims.Issuer).Val()
		if otpSecret == "" {
			fmt.Printf("[-] OTP: No secret found\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "OTP error",
			})
		}
		if !auth.VerifyOTP(otpSecret, passwordRequest.Otp) {
			// event password info failed otp
			err = event.NewEvent(clientPostgresDb, "Password info request failed", "User used an invalid OTP code", c.IP(), claims.Issuer)
			if err != nil {
				fmt.Printf("[!] Error occurred creating event: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}

			fmt.Printf("[-] OTP: Invalid code\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid OTP",
			})
		}

		// event password info success otp
		err = event.NewEvent(clientPostgresDb, "Password info request success", fmt.Sprintf("User viewed a password for %s protected by OTP", passwordRequest.Domain), c.IP(), claims.Issuer)
		if err != nil {
			fmt.Printf("[!] Error occurred creating event: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		fmt.Printf("[+] OTP: Access granted\n")
	}

	// event password info
	err = event.NewEvent(clientPostgresDb, "Password info request success", fmt.Sprintf("User viewed a password for %s", passwordRequest.Domain), c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	fmt.Printf("[+] Password info found\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Password found",
		"password": result,
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	var passwordRequest models.PasswordUpdate
	if err := c.BodyParser(&passwordRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
		})
	}

	// get user id from jwt
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Update password: Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	// get password from db
	result, err := clientPostgresDb.Password.FindMany(
		db.Password.ID.Equals(passwordRequest.PasswordId),
		db.Password.UserID.Equals(claims.Issuer),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Update password: Error occurred finding password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding password",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] Update password: No password found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No password found",
		})
	}

	// check if otp is needed
	if result[0].OtpProtected {
		fmt.Printf("[+] Update password: Password is protected by OTP\n")
		otpSecret := clientRedisDb.Get(context.Background(), claims.Issuer).Val()
		if otpSecret == "" {
			fmt.Printf("[-] Update password: OTP: No secret found\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "OTP error",
			})
		}
		
		if !auth.VerifyOTP(otpSecret, passwordRequest.Otp) {
			// event password update failed otp
			err = event.NewEvent(clientPostgresDb, "Password update failed", "User used an invalid OTP code", c.IP(), claims.Issuer)
			if err != nil {
				fmt.Printf("[!] Update password: Error occurred creating event: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}

			fmt.Printf("[-] Update password: OTP: Invalid code\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid OTP",
			})
		}

		fmt.Printf("[+] Update password: OTP: Access granted\n")
	}

	// check if the old password given is equal to the password in the db
	if !cryptography.CompareStrings(passwordRequest.OldPassword, result[0].Password) {
		fmt.Printf("[-] Update password: Password not match\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Password not match, old password is incorrect",
		})
	}

	// event password update
	err = event.NewEvent(clientPostgresDb, "Password updated", fmt.Sprintf("User updated a password for %s", result[0].Website), c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Update password: Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// update password
	_, err = clientPostgresDb.Password.FindMany(
		db.Password.ID.Equals(passwordRequest.PasswordId),
		db.Password.UserID.Equals(claims.Issuer),
	).Update(
		db.Password.Password.Set(passwordRequest.NewPassword),
		db.Password.OtpProtected.Set(passwordRequest.OtpProtected),
		db.Password.Description.Set(passwordRequest.NewDescription),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Update password: Error occurred updating password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating password",
		})
	}

	fmt.Printf("[+] Update password: Password updated successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

func DeletePassword(c *fiber.Ctx) error {
	var passwordRequest models.PasswordDelete
	if err := c.BodyParser(&passwordRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
		})
	}

	// get user id from jwt
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Update password: Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	// get password from db
	result, err := clientPostgresDb.Password.FindMany(
		db.Password.ID.Equals(passwordRequest.PasswordId),
		db.Password.UserID.Equals(claims.Issuer),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Delete password: Error occurred finding password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding password",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] Delete password: No password found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No password found",
		})
	}

	// check if otp is needed
	if result[0].OtpProtected {
		fmt.Printf("[+] Delete password: Password is protected by OTP\n")
		otpSecret := clientRedisDb.Get(context.Background(), claims.Issuer).Val()
		if otpSecret == "" {
			fmt.Printf("[-] Delete password: OTP: No secret found\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "OTP error",
			})
		}
		
		if !auth.VerifyOTP(otpSecret, passwordRequest.Otp) {
			// event password delete failed otp
			err = event.NewEvent(clientPostgresDb, "Password delete failed", "User used an invalid OTP code", c.IP(), claims.Issuer)
			if err != nil {
				fmt.Printf("[!] Delete password: Error occurred creating event: %s", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Internal server error",
				})
			}

			fmt.Printf("[-] Delete password: OTP: Invalid code\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid OTP",
			})
		}

		fmt.Printf("[+] Delete password: OTP: Access granted\n")
	}

	// event password delete
	err = event.NewEvent(clientPostgresDb, "Password deleted", fmt.Sprintf("User deleted a password for %s", result[0].Website), c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Delete password: Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// delete password
	_, err = clientPostgresDb.Password.FindMany(
		db.Password.ID.Equals(passwordRequest.PasswordId),
		db.Password.UserID.Equals(claims.Issuer),
	).Delete().Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Delete password: Error occurred deleting password: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting password",
		})
	}

	fmt.Printf("[+] Delete password: Password deleted successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password deleted successfully",
	})
}

func Logout(c *fiber.Ctx) error {
	// add the JWT token to redis blacklist
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	// check if cookie is null
	if cookie == "" {
		fmt.Printf("[+] Logout: cookie is null")
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Cookie is null",
		})
	}

	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	tokenRemainTime := auth.TokenRemainingTime(claims)
	// check if token is already expired
	if tokenRemainTime <= 0 {
		fmt.Printf("[+] Logout: token is already expired")
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "Token is already expired",
		})
	}

	_, err = clientRedisDb.Set(context.Background(), cookie, claims.Issuer, time.Duration(tokenRemainTime)*time.Second).Result()
	if err != nil {
		fmt.Printf("[!] Error occurred adding token to Redis blacklist: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	fmt.Printf("[+] Token added to Redis blacklist\n")

	// remove token from postgres db for sessions online
	_, err = clientPostgresDb.Token.FindMany(
		db.Token.TokenValue.Equals(cookie),
	).Delete().Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Session: Error occurred removing token from db: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// clear the cookie
	cookieClear := fiber.Cookie{
		Name:     os.Getenv("JWT_COOKIE_TOKEN_NAME"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookieClear)

	fmt.Printf("[+] User %s logout\n", claims.Issuer)
	fmt.Printf("[+] Session: Token removed from db\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Logout successful",
	})
}

func CheckPasswordLeak(c *fiber.Ctx) error {
	var passwordCheck models.PasswordLeakCheck

	if err := c.BodyParser(&passwordCheck); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(passwordCheck) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "missing required fields",
		})
	}

	// check if the password is in the password leaked database
	result, err := clientPostgresDb.PasswordLeak.FindMany(
		db.PasswordLeak.PasswordHash.Contains(passwordCheck.PasswordPartialHash),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Error occurred finding password in leaked database: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(result) > 0 {
		possiblePasswordHash := []string{}
		// check if the partial hash is the start of the password hash
		for _, password := range result {
			if strings.HasPrefix(password.PasswordHash, passwordCheck.PasswordPartialHash) {
				possiblePasswordHash = append(possiblePasswordHash, password.PasswordHash)
			}
		}

		if len(possiblePasswordHash) > 0 {
			fmt.Printf("[+] Password found in leaked database\n")
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Possible password found in leaked database",
				"result":  true,
				"hashes": possiblePasswordHash,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password safe",
		"result": false,
	})
}

func GetEvents(c *fiber.Ctx) error {
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	var eventRequest models.EventRequest
	if err := c.BodyParser(&eventRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(eventRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
		})
	}

	// get events from db
	startDateTime, err := time.Parse(time.RFC3339, eventRequest.StartDateTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid start date format",
		})
	}
	endDateTime, err := time.Parse(time.RFC3339, eventRequest.EndDateTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid end date format",
		})
	}

	result, err := clientPostgresDb.Event.FindMany(
		db.Event.UserID.Equals(claims.Issuer),
		db.Event.CreatedAt.After(startDateTime),
		db.Event.CreatedAt.Before(endDateTime),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Error occurred getting events: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] No events found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No events found",
		})
	}

	fmt.Printf("[+] Events found\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Events found",
		"events":  result,
	})
}

func GetSessions(c *fiber.Ctx) error {
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	result, err := clientPostgresDb.Token.FindMany(
		db.Token.UserID.Equals(claims.Issuer),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Error occurred getting sessions: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	
	if len(result) == 0 {
		fmt.Printf("[-] No sessions found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No sessions found",
		})
	}

	// search for the current user's sessions
	userSessions := []models.SessionModelResponse{}
	for index, session := range result {
		userSessions = append(userSessions, models.SessionModelResponse{
			DatabaseElemID: session.ID,
			LastUse:   	session.UpdatedAt.Format(time.RFC3339),
			IpAddress:  session.IPAddress,
			CreatedAt: 	session.CreatedAt.Format(time.RFC3339),
			CurrentUser: false,
		})
		if session.TokenValue == cookie {
			userSessions[index].CurrentUser = true
		}
	}


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sessions found",
		"sessions": userSessions,
	})
}

func ForceLogoutSession(c *fiber.Ctx) error {
	var sessionDeleteRequest models.SessionDeleteRequest
	if err := c.BodyParser(&sessionDeleteRequest); err != nil {
		return fiber.ErrBadRequest
	}

	if !utils.CheckAllFieldsHaveValue(sessionDeleteRequest) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Missing required fields",
		})
	}
		
	// get user id from jwt
	cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))
	claims, err := auth.ParseJWTToken(cookie)
	if err != nil {
		fmt.Printf("[!] Force remove session: Error occurred parsing JWT token: %s", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Jwt error",
		})
	}

	// otp
	otpSecret := clientRedisDb.Get(context.Background(), claims.Issuer).Val()
	if otpSecret == "" {
		fmt.Printf("[-] Force remove session: OTP: No secret found\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "OTP error",
		})
	}

	if !auth.VerifyOTP(otpSecret, sessionDeleteRequest.Otp) {
		// event session delete failed otp
		err = event.NewEvent(clientPostgresDb, "Session delete failed", "User used an invalid OTP code", c.IP(), claims.Issuer)
		if err != nil {
			fmt.Printf("[!] Force remove session: Error occurred creating event: %s", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		fmt.Printf("[-] Force remove session: OTP: Invalid code\n")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	// event session delete
	err = event.NewEvent(clientPostgresDb, "Session deleted", "User deleted a session", c.IP(), claims.Issuer)
	if err != nil {
		fmt.Printf("[!] Force remove session: Error occurred creating event: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	// remove session from postgres db
	// get
	result, err := clientPostgresDb.Token.FindMany(
		db.Token.ID.Equals(sessionDeleteRequest.DatabaseElemID),
	).Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Force remove session: Error occurred finding session: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error finding session",
		})
	}

	if len(result) == 0 {
		fmt.Printf("[-] Force remove session: No session found\n")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No session found",
		})
	}

	// delete 
	_, err = clientPostgresDb.Token.FindMany(
		db.Token.ID.Equals(sessionDeleteRequest.DatabaseElemID),
	).Delete().Exec(context.Background())
	if err != nil {
		fmt.Printf("[!] Force remove session: Error occurred deleting session: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting session",
		})
	}

	// remove session from redis db
	_, err = clientRedisDb.Del(context.Background(), result[0].TokenValue).Result()
	if err != nil {
		fmt.Printf("[!] Force remove session: Error occurred deleting session from Redis: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}

	fmt.Printf("[+] Force remove session: Session removed successfully\n")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Session Removed successfully",
	})
}