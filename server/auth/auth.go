package auth

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/dipithedipi/password-manager/utils"
	"github.com/golang-jwt/jwt"
	"github.com/xlzd/gotp"
)


func GenerateTOTPSecret() string {
	secretLenght, err := strconv.Atoi(os.Getenv("OTP_SECRET_LENGTH"))
	if err != nil {panic(err)}
	return gotp.RandomSecret(secretLenght)
}

func GenerateUriTOTPWithSecret(randomSecret string, email string) string {
	return gotp.NewDefaultTOTP(randomSecret).ProvisioningUri(email, "Password Manager")
}
   
func VerifyOTP(randomSecret string, userTotp string) bool {
	totp := gotp.NewDefaultTOTP(randomSecret)

	// Validate the provided OTP
	if totp.Verify(userTotp, time.Now().Unix()) {
		return true
	} else {
		return false
	}
}

// JWT
func GenerateJWTToken(userId string, validPeriod string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userId,
		ExpiresAt: utils.CalculateExpireTimeInt64(validPeriod),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseJWTToken(token string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, err
	}
	return claims, nil
}

func VerifyJWTToken(token string) (bool, error) {
	claims, err := ParseJWTToken(token)
	if err != nil {
		return false, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return false, nil
	}
	return true, nil
}

func TokenRemainingTime(token *jwt.StandardClaims) int64 {
	return token.ExpiresAt - time.Now().Unix()
}

func MiddlewareJWTAuth(clientRedisDb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Print("[+] Middleware JWT Auth\n")	

		// Ottieni il token dalla cookie
		cookie := c.Cookies(os.Getenv("JWT_COOKIE_TOKEN_NAME"))

		// check if the cookie is empty
		if cookie == "" {
			fmt.Print("[!] Middleware JWT Auth: token is missing\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is missing",
			})
		}

		// check if the token is valid
		valid, err := VerifyJWTToken(cookie)
		if err != nil {
			fmt.Printf("[!] Middleware JWT Auth: token is invalid\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is invalid",
			})
		}
		if !valid {
			fmt.Print("[!] Middleware JWT Auth: token is expired\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is expired",
			})
		}

		// check if the token is in the blacklist
		result, err := clientRedisDb.Get(context.Background(), cookie).Result()
		if err == redis.Nil {
			// token not found in the blacklist, continue
			fmt.Print("[+] Middleware JWT Auth: sucess\n")
			return c.Next()
		} else if err != nil {
			fmt.Printf("[!] Error check redis blacklist: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal server error",
			})
		}

		// token found in the blacklist
		if result != "" {
			fmt.Print("[!] Middleware JWT Auth: token found in the blacklist\n")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token is expired",
			})
		}

		// unexpected error
		fmt.Print("[!] Middleware JWT Auth: unexpected error\n")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
}
