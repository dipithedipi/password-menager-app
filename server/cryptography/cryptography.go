package cryptography

import (
	"github.com/dipithedipi/password-manager/models"
    "github.com/dipithedipi/password-manager/cryptography/keys"
	"golang.org/x/crypto/argon2"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "errors"
    "crypto/cipher"
    "bytes"
    "encoding/base64"
    "crypto/aes"
    "encoding/pem"
    "os"
	"fmt"
)

// RSA
func GenerateKeysRSA(publicKeyPath string, privateKeyPath string, keyLenght int) {
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

func EncryptDataRSA(plaintext []byte) []byte {
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

func DecryptDataRSA(ciphertext []byte, privateKeyPEM []byte) []byte {
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

// AES CBC
func EncryptAESCBC(plaintext string, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("key error1", err)
        return nil, err
	}

	if plaintext == "" {
		fmt.Println("plain content empty")
        return nil, errors.New("plain content empty")
	}

    initialVector := key[:aes.BlockSize]
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(string(plaintext))
	content = PKCS5Padding(content, block.BlockSize())
	ciphertext := make([]byte, len(content))
	ecb.CryptBlocks(ciphertext, content)

	return ciphertext, nil
}

func DecryptAESCBC(ciphertext []byte, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
	if err != nil {
        return nil, err
	}
	if len(ciphertext) == 0 {
        return nil, errors.New("plain content empty")
	}

    initialVector := key[:aes.BlockSize]
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(ciphertext))
	ecb.CryptBlocks(decrypted, ciphertext)

	return PKCS5Trimming(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

// Base64
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) []byte {
    decoded, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        panic(err)
    }

    return decoded
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