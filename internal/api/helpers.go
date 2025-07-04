package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

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

func (a *App) LogTemplateError(handler string, template string, err error) {
	a.Logger.Error("template execution error",
		"handler", handler,
		"template", template,
		"error", err,
	)
}

func (a *App) LogInternalServerError(r *http.Request, msg string, handler string, err error) {

	a.Logger.Error(msg,
		"handler", handler,
		"path", r.URL.Path,
		"status", 500,
		"method", r.Method,
		"error", err,
	)
}

func (a *App) LogBadRequest(r *http.Request, msg string, handler string, err error) {
	a.Logger.Error(msg,
		"handler", handler,
		"path", r.URL.Path,
		"status", 400,
		"method", r.Method,
		"error", err,
	)
}
