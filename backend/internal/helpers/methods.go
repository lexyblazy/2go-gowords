package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func NewUUIDV4() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant bits (RFC 4122)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}

func GenerateRecoveryCode() (string, error) {
	b := make([]byte, 12)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	code := string(b)

	return fmt.Sprintf("%s-%s-%s",
		code[0:4],
		code[4:8],
		code[8:12],
	), nil
}

func HashRecoveryCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func ValidateRecoveryCode(inputCode string, storedHash string) bool {
	inputHash := HashRecoveryCode(inputCode)
	return inputHash == storedHash
}
