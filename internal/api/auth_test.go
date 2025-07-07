package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func createTestKeyFiles(t *testing.T) (string, string) {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test key: %v", err)
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	privateKeyFile, err := os.CreateTemp("", "private_key_*.pem")
	if err != nil {
		t.Fatalf("failed to create temp private key file: %v", err)
	}

	publicKeyFile, err := os.CreateTemp("", "public_key_*.pem")
	if err != nil {
		t.Fatalf("failed to create temp public key file: %v", err)
	}

	if _, err := privateKeyFile.Write(privateKeyPEM); err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}
	if _, err := publicKeyFile.Write(publicKeyPEM); err != nil {
		t.Fatalf("failed to write public key: %v", err)
	}

	privateKeyFile.Close()
	publicKeyFile.Close()

	return privateKeyFile.Name(), publicKeyFile.Name()
}

func TestInitKeys(t *testing.T) {
	tests := []struct {
		name         string
		setupKeyFile func(t *testing.T) string
		expectError  bool
	}{
		{
			name: "valid private key file",
			setupKeyFile: func(t *testing.T) string {
				privateKeyPath, _ := createTestKeyFiles(t)
				return privateKeyPath
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPrivateKey := privateKey
			privateKey = nil
			defer func() { privateKey = originalPrivateKey }()

			keyFile := tt.setupKeyFile(t)
			defer os.Remove(keyFile)

			originalWD, _ := os.Getwd()
			defer os.Chdir(originalWD)

			defer os.Remove("./private_key.pem")
			err := os.Link(keyFile, "./private_key.pem")
			if err != nil {
				data, readErr := os.ReadFile(keyFile)
				if readErr != nil {
					t.Fatalf("failed to read key file: %v", readErr)
				}
				writeErr := os.WriteFile("./private_key.pem", data, 0600)
				if writeErr != nil {
					t.Fatalf("failed to write key file: %v", writeErr)
				}
			}

			InitKeys()

			if privateKey == nil {
				t.Error("expected privateKey to be set but it was nil")
			}
		})
	}
}

func TestIssueKeys(t *testing.T) {
	privateKeyPath, _ := createTestKeyFiles(t)
	defer os.Remove(privateKeyPath)

	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("failed to read test private key: %v", err)
	}

	testPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		t.Fatalf("failed to parse test private key: %v", err)
	}

	tests := []struct {
		name        string
		setupKey    func()
		expectError bool
	}{
		{
			name: "valid private key",
			setupKey: func() {
				privateKey = testPrivateKey
			},
			expectError: false,
		},
		{
			name: "nil private key",
			setupKey: func() {
				privateKey = nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalPrivateKey := privateKey
			defer func() { privateKey = originalPrivateKey }()

			tt.setupKey()

			var token string
			var err error

			if tt.expectError && tt.name == "nil private key" {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic but didn't get one")
					}
				}()
			}

			token, err = IssueKeys()

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if token != "" {
					t.Errorf("expected empty token but got %s", token)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token == "" {
					t.Error("expected token but got empty string")
				}

				parsedToken, parseErr := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
					return &testPrivateKey.PublicKey, nil
				})

				if parseErr != nil {
					t.Errorf("failed to parse issued token: %v", parseErr)
				}

				if !parsedToken.Valid {
					t.Error("issued token is not valid")
				}

				claims, ok := parsedToken.Claims.(*JwtClaims)
				if !ok {
					t.Error("failed to extract claims from token")
				} else {
					if claims.Access_level != 4 {
						t.Errorf("expected access level 4 but got %d", claims.Access_level)
					}
					if claims.Subject != "authentication" {
						t.Errorf("expected subject 'authentication' but got '%s'", claims.Subject)
					}
					if claims.ExpiresAt == nil {
						t.Error("expected expiration time but got nil")
					} else {
						expectedExpiry := time.Now().Add(time.Hour)
						actualExpiry := claims.ExpiresAt.Time
						timeDiff := actualExpiry.Sub(expectedExpiry)
						if timeDiff < -time.Minute || timeDiff > time.Minute {
							t.Errorf("expected expiry around %v but got %v", expectedExpiry, actualExpiry)
						}
					}
				}
			}
		})
	}
}

func TestAuthJwtClaims(t *testing.T) {
	tests := []struct {
		name        string
		accessLevel int
		subject     string
		expiry      time.Time
	}{
		{
			name:        "standard claims",
			accessLevel: 4,
			subject:     "authentication",
			expiry:      time.Now().Add(time.Hour),
		},
		{
			name:        "admin claims",
			accessLevel: 0,
			subject:     "admin",
			expiry:      time.Now().Add(2 * time.Hour),
		},
		{
			name:        "user claims",
			accessLevel: 3,
			subject:     "user",
			expiry:      time.Now().Add(30 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := &JwtClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   tt.subject,
					ExpiresAt: jwt.NewNumericDate(tt.expiry),
				},
				Access_level: tt.accessLevel,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			if token == nil {
				t.Error("expected JWT token to be created but got nil")
			}

			if claims.Access_level != tt.accessLevel {
				t.Errorf("expected access level %d but got %d", tt.accessLevel, claims.Access_level)
			}
			if claims.Subject != tt.subject {
				t.Errorf("expected subject '%s' but got '%s'", tt.subject, claims.Subject)
			}
			if claims.ExpiresAt == nil {
				t.Error("expected expiration time but got nil")
			}
		})
	}
}

func TestAuthJwtClaimsJSON(t *testing.T) {
	claims := &JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "test",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Access_level: 2,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	parsedClaims, ok := parsedToken.Claims.(*JwtClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if parsedClaims.Access_level != claims.Access_level {
		t.Errorf("expected access level %d but got %d", claims.Access_level, parsedClaims.Access_level)
	}
	if parsedClaims.Subject != claims.Subject {
		t.Errorf("expected subject '%s' but got '%s'", claims.Subject, parsedClaims.Subject)
	}
}

// Test edge cases and integration scenarios
func TestAuthIntegration(t *testing.T) {
	privateKeyPath, publicKeyPath := createTestKeyFiles(t)
	defer os.Remove(privateKeyPath)
	defer os.Remove(publicKeyPath)

	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("failed to read private key: %v", err)
	}

	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		t.Fatalf("failed to read public key: %v", err)
	}

	err = os.WriteFile("./private_key.pem", privateKeyData, 0600)
	if err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}
	defer os.Remove("./private_key.pem")

	err = os.WriteFile("./public_key.pem", publicKeyData, 0644)
	if err != nil {
		t.Fatalf("failed to write public key: %v", err)
	}
	defer os.Remove("./public_key.pem")

	originalPrivateKey := privateKey
	defer func() { privateKey = originalPrivateKey }()

	InitKeys()

	token, err := IssueKeys()
	if err != nil {
		t.Fatalf("failed to issue token: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		t.Fatalf("failed to parse public key: %v", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("token should be valid but isn't")
	}

	claims, ok := parsedToken.Claims.(*JwtClaims)
	if !ok {
		t.Fatal("failed to extract claims")
	}

	if claims.Access_level != 4 {
		t.Errorf("expected access level 4 but got %d", claims.Access_level)
	}
}
