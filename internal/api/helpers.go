package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func generateRandomString() (string, error) {
	randomString := make([]byte, 8)
	_, err := rand.Read(randomString)
	if err != nil {
		return "", fmt.Errorf("failed to generate a string: %+v", err)
	}
	encodedString := base64.URLEncoding.EncodeToString(randomString)
	return string(encodedString), nil
}

func hashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(hash), err

}

func verifyPasswordHash(password string, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
