package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/dipithedipi/password-manager/cryptography/keys"
	"github.com/dipithedipi/password-manager/models"
	"golang.org/x/crypto/argon2"
)

// RSA OEAP
// private key: PKCS1
// public key: SubjectPublicKeyInfo structure (PKIX)
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

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey.(*rsa.PublicKey), plaintext, nil)
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

    plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
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

// Sha1
func Sha1(data string) string {
    h := crypto.SHA1.New()
    h.Write([]byte(data))
    return fmt.Sprintf("%x", h.Sum(nil))
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

func CompareStrings(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}