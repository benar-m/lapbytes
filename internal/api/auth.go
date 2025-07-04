package api

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	Access_level int `json:"access_level,omitempty"`
}

var (
	privateKey *rsa.PrivateKey
)

func InitKeys() {
	privateKeyData, err := os.ReadFile("./private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		log.Fatal(err)
	}
}

func IssueKeys() (jwtToken string, err error) {

	claims := &JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "authentication",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Access_level: 4, //work to do
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
