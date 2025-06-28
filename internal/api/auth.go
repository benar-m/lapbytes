package api

import (
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Access_level int `json:"access_level,omitempty"`
}

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

func InitKeys() {
	publicKeyData, err := os.ReadFile("./public_key.pem")
	if err != nil {
		log.Fatal(err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		log.Fatal(err)
	}

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
	//Only called after a succesful Authentication -- To further Mod to switch key roles based on access_level

	claims := &jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "authentication",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Access_level: 1, //work to do
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return t, nil
}
