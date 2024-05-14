package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken(tokenLength int) (string, error) {
	rawLength := 3 * tokenLength / 4

	randomBytes := make([]byte, rawLength)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
