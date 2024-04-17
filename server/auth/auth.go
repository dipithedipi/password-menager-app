package auth

import (
	"fmt"
	"time"
	"os"
	"strconv"
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
		fmt.Println("[+] OTP: Access granted")
		return true
	} else {
		fmt.Println("[-] OTP: Access denied")
		return false
}
}