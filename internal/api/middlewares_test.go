package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setupTestAppForMiddleware() *App {
	handler := slog.NewTextHandler(os.Stderr, nil)
	logger := slog.New(handler)

	return &App{
		DB:     nil,
		Logger: logger,
	}
}

func generateTestKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func createTestToken(privateKey *rsa.PrivateKey, accessLevel int) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Access_level: accessLevel,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func TestReqLoggingMW(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET request",
			method: "GET",
			path:   "/test",
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "/api/users",
		},
		{
			name:   "PUT request",
			method: "PUT",
			path:   "/api/products/1",
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			path:   "/api/users/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{}

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			handler := app.ReqLoggingMW(testHandler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d but got %d", http.StatusOK, w.Code)
			}

			expectedBody := "test response"
			if w.Body.String() != expectedBody {
				t.Errorf("expected body '%s' but got '%s'", expectedBody, w.Body.String())
			}
		})
	}
}

func TestGeneralJwtVerifierMW(t *testing.T) {
	privateKey, publicKey, err := generateTestKeys()
	if err != nil {
		t.Fatalf("failed to generate test keys: %v", err)
	}

	originalPublicKey := PublicKey
	PublicKey = publicKey
	defer func() { PublicKey = originalPublicKey }()

	tests := []struct {
		name               string
		authHeader         string
		setupToken         func() string
		expectedStatus     int
		expectedError      string
		expectContextValue bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			setupToken:     func() string { return "" },
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "missing authorization header",
		},
		{
			name:           "invalid authorization header format",
			authHeader:     "InvalidFormat",
			setupToken:     func() string { return "" },
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid authorization header",
		},
		{
			name:           "invalid bearer format",
			authHeader:     "Basic token123",
			setupToken:     func() string { return "" },
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid authorization header",
		},
		{
			name:       "valid token",
			authHeader: "Bearer %s",
			setupToken: func() string {
				token, err := createTestToken(privateKey, 2)
				if err != nil {
					t.Fatalf("failed to create test token: %v", err)
				}
				return token
			},
			expectedStatus:     http.StatusOK,
			expectContextValue: true,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid.token.here",
			setupToken:     func() string { return "" },
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "not authorized",
		},
		{
			name:       "expired token",
			authHeader: "Bearer %s",
			setupToken: func() string {
				claims := jwtClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
					Access_level: 2,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				tokenString, err := token.SignedString(privateKey)
				if err != nil {
					t.Fatalf("failed to create expired token: %v", err)
				}
				return tokenString
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "not authorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestAppForMiddleware()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectContextValue {
					claims := r.Context().Value(jwtClaimsKey)
					if claims == nil {
						t.Error("expected jwt claims in context but got nil")
					}
					if jwtClaims, ok := claims.(*jwtClaims); !ok {
						t.Error("expected jwt claims to be of type *jwtClaims")
					} else if jwtClaims.Access_level != 2 {
						t.Errorf("expected access level 2 but got %d", jwtClaims.Access_level)
					}
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			handler := app.GeneralJwtVerifierMW(testHandler)

			var authHeaderValue string
			if tt.authHeader != "" {
				token := tt.setupToken()
				if strings.Contains(tt.authHeader, "%s") {
					authHeaderValue = strings.Replace(tt.authHeader, "%s", token, 1)
				} else {
					authHeaderValue = tt.authHeader
				}
			}

			req := httptest.NewRequest("GET", "/test", nil)
			if authHeaderValue != "" {
				req.Header.Set("Authorization", authHeaderValue)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if response["error"] != tt.expectedError {
					t.Errorf("expected error '%s' but got '%s'", tt.expectedError, response["error"])
				}

				expectedContentType := "application/json"
				if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
					t.Errorf("expected content type '%s' but got '%s'", expectedContentType, contentType)
				}
			}
		})
	}
}

func TestIsAdminJwtVerifierMW(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid admin access (level 0)",
			setupContext: func() context.Context {
				claims := &jwtClaims{Access_level: 0}
				return context.WithValue(context.Background(), jwtClaimsKey, claims)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid admin access (level 1)",
			setupContext: func() context.Context {
				claims := &jwtClaims{Access_level: 1}
				return context.WithValue(context.Background(), jwtClaimsKey, claims)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid access level (level 2)",
			setupContext: func() context.Context {
				claims := &jwtClaims{Access_level: 2}
				return context.WithValue(context.Background(), jwtClaimsKey, claims)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
		{
			name: "invalid access level (level 4)",
			setupContext: func() context.Context {
				claims := &jwtClaims{Access_level: 4}
				return context.WithValue(context.Background(), jwtClaimsKey, claims)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
		{
			name: "missing jwt claims in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
		{
			name: "wrong context key",
			setupContext: func() context.Context {
				claims := &jwtClaims{Access_level: 0}
				return context.WithValue(context.Background(), "wrong_key", claims)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
		{
			name: "nil claims",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), jwtClaimsKey, nil)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
		{
			name: "wrong type in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), jwtClaimsKey, "not_jwt_claims")
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "admin access required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := setupTestAppForMiddleware()

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("admin access granted"))
			})

			handler := app.IsAdminJwtVerifierMW(testHandler)

			req := httptest.NewRequest("GET", "/admin", nil)
			req = req.WithContext(tt.setupContext())
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if response["error"] != tt.expectedError {
					t.Errorf("expected error '%s' but got '%s'", tt.expectedError, response["error"])
				}

				expectedContentType := "application/json"
				if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
					t.Errorf("expected content type '%s' but got '%s'", expectedContentType, contentType)
				}
			} else {
				expectedBody := "admin access granted"
				if w.Body.String() != expectedBody {
					t.Errorf("expected body '%s' but got '%s'", expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestContextKey(t *testing.T) {
	ctx := context.WithValue(context.Background(), jwtClaimsKey, "test_value")
	value := ctx.Value(jwtClaimsKey)

	if value != "test_value" {
		t.Errorf("expected 'test_value' but got %v", value)
	}

	ctx2 := context.WithValue(context.Background(), contextKey("jwt_claims"), "test_value2")
	value2 := ctx2.Value(jwtClaimsKey)

	if value2 != "test_value2" {
		t.Errorf("expected 'test_value2' but got %v", value2)
	}
}

func TestJwtClaims(t *testing.T) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Access_level: 1,
	}

	if claims.Access_level != 1 {
		t.Errorf("expected access level 1 but got %d", claims.Access_level)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token == nil {
		t.Error("expected JWT token to be created but got nil")
	}
}
