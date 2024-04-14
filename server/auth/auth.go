package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/dipithedipi/password-manager/models"
	"github.com/xlzd/gotp"
	"golang.org/x/crypto/argon2"
)

func GenerateFromPassword(password string, p *models.ArgonParams) (hash []byte, salt []byte, err error) {
	// Generate a cryptographically secure random salt.
	salt, err = GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return nil, nil, err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash = argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	return hash, salt, nil
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Salt(lenght uint32) []byte {
	salt, err := GenerateRandomBytes(lenght)
	if err != nil {
		fmt.Printf("[!] ERROR: generating salt %v", err)
		panic(err)
	}

	return salt
}

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