package keys

import (
	"os"
)

func ReadPublicKeyPEM() []byte {
    publicKeyPEM, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
    if err != nil {
        panic(err)
    }

    return publicKeyPEM
}

func ReadPrivateKeyPEM() []byte {
    publicKeyPEM, err := os.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
    if err != nil {
        panic(err)
    }

    return publicKeyPEM
}