package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() (string, error) {
    salt := make([]byte, 16)
    _, err := rand.Read(salt)
    if err != nil {
        return "", err
    }
    return base64.RawStdEncoding.EncodeToString(salt), nil
}

func HashPassword(password string) string {
	var salt string = os.Getenv("PASSWORD_SALT")
    hash := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
    return base64.RawStdEncoding.EncodeToString(hash)
}

func CheckPassword(password, hashedPassword string) (bool) {
    computedHash := HashPassword(password)
    return computedHash == hashedPassword
}