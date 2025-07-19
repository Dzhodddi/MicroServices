package security

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateToken() string {
	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])
	return hashedToken
}
