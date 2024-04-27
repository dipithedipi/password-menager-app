package cryptography

import (
	"github.com/dipithedipi/password-manager/models"
    "github.com/dipithedipi/password-manager/cryptography/keys"
	"golang.org/x/crypto/argon2"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "errors"
    "strings"
    "crypto/cipher"
    "bytes"
    "crypto/subtle"
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

func EncryptServerDataRSA(plaintext []byte) ([]byte, error) {
    publicKeyPEM := keys.ReadPublicKeyPEM()
    return EncryptDataRSA(plaintext, publicKeyPEM)
}

func DecryptServerDataRSA(ciphertext []byte) []byte {
    privateKeyPEM := keys.ReadPrivateKeyPEM()
    return DecryptDataRSA(ciphertext, privateKeyPEM)
}

func EncryptDataRSA(plaintext []byte, publicKeyPEM []byte) ([]byte, error) {
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
        return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), plaintext)
	if err != nil {
        return nil, err
	}

	return ciphertext, nil
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

func ConvertBase64PublicKeyToPEM(base64Key string) (string, error) {
    // Decode the base64-encoded public key
    decodedKey, err := base64.StdEncoding.DecodeString(base64Key)
    if err != nil {
        return "", fmt.Errorf("failed to decode public key: %v", err)
    }

    // add header and footer to the key
    pemKey := string(decodedKey[:])
    pemKey = "-----BEGIN RSA PUBLIC KEY-----\n" + pemKey + "\n-----END RSA PUBLIC KEY-----"

    return pemKey, nil
}

// AES CBC (not working)
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

// AES GCM
func EncryptAESGCM(plaintext string, secretKey []byte) ([]byte, error)  {
    fmt.Println("plaintext lenght: %v", len(plaintext))
    fmt.Println("plaintext: %v", plaintext)
    fmt.Println("key lenght: %v", len(secretKey))
    fmt.Println("key: %v", secretKey)

    aes, err := aes.NewCipher(secretKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return nil, err
    }

    // We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
    // A nonce should always be randomly generated for every encryption.
    nonce := make([]byte, gcm.NonceSize())
    _, err = rand.Read(nonce)
    if err != nil {
        return nil, err
    }

    // ciphertext here is actually nonce+ciphertext
    // So that when we decrypt, just knowing the nonce size
    // is enough to separate it from the ciphertext.
    ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

    return ciphertext, nil
}

func DecryptAESGCM(ciphertext string, secretKey []byte) ([]byte, error) {
    aes, err := aes.NewCipher(secretKey)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(aes)
    if err != nil {
        return nil, err
    }

    // Since we know the ciphertext is actually nonce+ciphertext
    // And len(nonce) == NonceSize(). We can separate the two.
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
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

// Argon2
func HashPassword(password string, p *models.ArgonParams) (encodedHash string, err error) {
    salt, err := generateRandomBytes(p.SaltLength)
    if err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

    // Base64 encode the salt and hashed password.
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    // Return a string using the standard encoded hash representation.
    encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

    return encodedHash, nil
}

func ComparePasswordAndHash(password string, encodedHash string) (match bool, err error) {
    // Extract the parameters, salt and derived key from the encoded password
    // hash.
    p, salt, hash, err := DecodeHash(encodedHash)

    if err != nil {
        return false, err
    }

    // Derive the key from the other password using the same parameters.
    otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

    // Check that the contents of the hashed passwords are identical. Note
    // that we are using the subtle.ConstantTimeCompare() function for this
    // to help prevent timing attacks.
    if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
        return true, nil
    }
    return false, nil
}

func DecodeHash(encodedHash string) (p *models.ArgonParams, salt, hash []byte, err error) {
    vals := strings.Split(encodedHash, "$")
    if len(vals) != 6 {
        return nil, nil, nil, errors.New("the encoded hash is not in the correct format")
    }

    var version int
    _, err = fmt.Sscanf(vals[2], "v=%d", &version)
    if err != nil {
        return nil, nil, nil, err
    }
    if version != argon2.Version {
        return nil, nil, nil, errors.New("incompatible version of argon2")
    }

    p = &models.ArgonParams{}
    _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
    if err != nil {
        return nil, nil, nil, err
    }

    salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
    if err != nil {
        return nil, nil, nil, err
    }
    p.SaltLength = uint32(len(salt))

    hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
    if err != nil {
        return nil, nil, nil, err
    }
    p.KeyLength = uint32(len(hash))

    return p, salt, hash, nil
}

// other
func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateSalt(lenght uint32) ([]byte, error) {
	salt, err := generateRandomBytes(lenght)
	if err != nil {
		fmt.Printf("[!] ERROR: generating salt %v", err)
		return nil, err
	}

	return salt, nil
}