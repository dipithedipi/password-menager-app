package cryptography

import (
	"github.com/dipithedipi/password-manager/models"
    "github.com/dipithedipi/password-manager/cryptography/keys"
	"golang.org/x/crypto/argon2"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "os"
	"fmt"
)

func GenerateKeys(publicKeyPath string, privateKeyPath string, keyLenght int) {
    privateKey, err := rsa.GenerateKey(rand.Reader, keyLenght)
    if err != nil {
        panic(err)
    }

    publicKey := &privateKey.PublicKey

    privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privateKeyBytes,
    })
    err = os.WriteFile(privateKeyPath, privateKeyPEM, 0644)
    if err != nil {
        panic(err)
    }

    publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        panic(err)
    }
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    err = os.WriteFile(publicKeyPath, publicKeyPEM, 0644)
    if err != nil {
        panic(err)
    }
}

func EncryptData(plaintext []byte) []byte {
    publicKeyPEM := keys.ReadPublicKeyPEM()
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
		panic(err)
	}

	return ciphertext
}

func DecryptData(ciphertext []byte, privateKeyPEM []byte) []byte {
    privateKeyBlock, _ := pem.Decode(privateKeyPEM)
    privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
    if err != nil {
        panic(err)
    }

    plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
    if err != nil {
        panic(err)
    }

    return plaintext
}

func HashPassword(password string, p *models.ArgonParams) (hash []byte, salt []byte, err error) {
	// Generate a cryptographically secure random salt.
	salt, err = generateRandomBytes(p.SaltLength)
	if err != nil {
		return nil, nil, err
	}

	// Pass the plaintext password, salt and parameters to the argon2.IDKey
	// function. This will generate a hash of the password using the Argon2id
	// variant.
	hash = argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	return hash, salt, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateSalt(lenght uint32) []byte {
	salt, err := generateRandomBytes(lenght)
	if err != nil {
		fmt.Printf("[!] ERROR: generating salt %v", err)
		panic(err)
	}

	return salt
}