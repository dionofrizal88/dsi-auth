package security

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

// HashPasswordWithSalt hashes a password with a salt using SHA-256
func HashPasswordWithSalt(password, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateRecoveryToken hashes a token with a salt using SHA-256
func GenerateRecoveryToken(str, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(str + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateSalt creates a random salt of the specified length
func GenerateSalt(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(salt)
}
