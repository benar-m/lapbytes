package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateRandomString() (string, error) {
	randomString := make([]byte, 32)
	_, err := rand.Read(randomString)
	if err != nil {
		return "", fmt.Errorf("failed to generate a string: %+v", err)
	}
	encodedString := base64.URLEncoding.EncodeToString(randomString)
	return string(encodedString), nil
}
