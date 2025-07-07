package api

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Access_level int `json:"access_level,omitempty"`
}
type contextKey string

const jwtClaimsKey contextKey = "jwt_claims"

// Public, potential error here
var PublicKey *rsa.PublicKey

func init() {
	pkData, err := os.ReadFile("./public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(pkData)
	if err != nil {
		log.Fatal(err)
	}

}

// ReqLoggingMW logs HTTP requests with method, path, and duration
func (a *App) ReqLoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))

	})
}

// GeneralJwtVerifierMW verifies JWT tokens and extracts claims for authenticated requests
func (a *App) GeneralJwtVerifierMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "missing authorization header",
			})
			return
		}
		splitToken := strings.Split(authorizationHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid authorization header",
			})
			return
		}
		tokenString := splitToken[1]
		claims := jwtClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			return PublicKey, nil
		})
		if err != nil || !token.Valid {
			a.Logger.Error("jwt verfication error",
				"time", time.Now(),
				"ip", r.RemoteAddr,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "not authorized",
			})
			return
		}
		ctx := context.WithValue(r.Context(), jwtClaimsKey, &claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IsAdminJwtVerifierMW ensures the authenticated user has admin privileges
func (a *App) IsAdminJwtVerifierMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(jwtClaimsKey)
		c, ok := v.(*jwtClaims)

		if !ok || c == nil || c.Access_level > 1 {
			a.Logger.Error("admin access denied",
				"time", time.Now(),
				"ip", r.RemoteAddr,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "admin access required",
			})
			return

		}

		next.ServeHTTP(w, r)
	})
}
