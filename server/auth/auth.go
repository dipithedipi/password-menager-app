package auth

import (
	"fmt"
	"time"
	"github.com/xlzd/gotp"
)


func GenerateTOTPSecret(secretLenght int) string {
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