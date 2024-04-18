package auth

import (
	"os"
	"strconv"
	"time"

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