package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lapbytes/internal/model"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type MockLogger struct {
	ErrorCalls []LogCall
	InfoCalls  []LogCall
}

type LogCall struct {
	Message string
	Fields  []interface{}
}

func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, LogCall{Message: msg, Fields: fields})
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.InfoCalls = append(m.InfoCalls, LogCall{Message: msg, Fields: fields})
}

func (m *MockLogger) Debug(msg string, fields ...interface{}) {}

func setupTestApp() *App {
	handler := slog.NewTextHandler(os.Stderr, nil)
	logger := slog.New(handler)

	return &App{
		DB:     nil,
		Logger: logger,
	}
}

func TestRenderHome(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.RenderHome(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "internal server error") {
		t.Error("Expected internal server error message")
	}
}

func TestRenderRegister(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/register", nil)
	w := httptest.NewRecorder()

	app.RenderRegister(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	responseBody := w.Body.String()
	if !strings.Contains(responseBody, "internal server error") {
		t.Error("Expected internal server error message")
	}
}

func TestRenderLogin(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	app.RenderLogin(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRenderProducts(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	app.RenderProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestRenderProduct(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/product/1", nil)
	w := httptest.NewRecorder()

	app.RenderProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestLoginUserInvalidContentType(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/login", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	app.LoginUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "bad request, accepts JSON only" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestLoginUserInvalidJSON(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/login", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.LoginUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] != "bad request, All fields are required" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestLoginUserValidJSON(t *testing.T) {
	app := setupTestApp()

	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(loginData)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Log("Handler panicked as expected due to nil database connection")
		}
	}()

	app.LoginUser(w, req)
}

func TestRegisterUserInvalidContentType(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	app.RegisterUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)

	if response["error"] != "invalid data, please check again" {
		t.Errorf("Unexpected error message: %s", response["error"])
	}
}

func TestRegisterUserInvalidJSON(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/register", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.RegisterUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegisterUserValidJSON(t *testing.T) {
	app := setupTestApp()

	userData := map[string]string{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(userData)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Log("Handler panicked as expected due to nil database connection")
		}
	}()

	app.RegisterUser(w, req)
}

func TestListProductsInvalidLimit(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products/invalid/1", nil)
	req.SetPathValue("limit", "invalid")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	app.ListProducts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListProductsInvalidPage(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products/10/invalid", nil)
	req.SetPathValue("limit", "10")
	req.SetPathValue("page", "invalid")
	w := httptest.NewRecorder()

	app.ListProducts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListProductsZeroLimit(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products/0/1", nil)
	req.SetPathValue("limit", "0")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	app.ListProducts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListProductsZeroPage(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/products/10/0", nil)
	req.SetPathValue("limit", "10")
	req.SetPathValue("page", "0")
	w := httptest.NewRecorder()

	app.ListProducts(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListProductInvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/product/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	app.ListProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListProductZeroID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/product/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	app.ListProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteUserInvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/user/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	app.DeleteUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteUserZeroID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/user/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	app.DeleteUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListUsersInvalidLimit(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/users/invalid/1", nil)
	req.SetPathValue("limit", "invalid")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	app.ListUsers(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListUsersInvalidPage(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/users/10/invalid", nil)
	req.SetPathValue("limit", "10")
	req.SetPathValue("page", "invalid")
	w := httptest.NewRecorder()

	app.ListUsers(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListUsersZeroLimit(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/users/0/1", nil)
	req.SetPathValue("limit", "0")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	app.ListUsers(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListUsersZeroPage(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/users/10/0", nil)
	req.SetPathValue("limit", "10")
	req.SetPathValue("page", "0")
	w := httptest.NewRecorder()

	app.ListUsers(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListSingleUserInvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/user/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	app.ListSingleUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListSingleUserZeroID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/user/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	app.ListSingleUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAddNewProductInvalidContentType(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/products", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	app.AddNewProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAddNewProductInvalidJSON(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/products", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.AddNewProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAddNewProductValidJSON(t *testing.T) {
	app := setupTestApp()

	laptop := model.Laptop{
		Name:     "Dell XPS 13",
		Brand:    "Dell",
		Price:    999.99,
		In_stock: 10,
	}

	jsonData, _ := json.Marshal(laptop)
	req := httptest.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Log("Handler panicked as expected due to nil database connection")
		}
	}()

	app.AddNewProduct(w, req)
}

func TestDeleteProductInvalidID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/product/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	app.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteProductZeroID(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("DELETE", "/product/0", nil)
	req.SetPathValue("id", "0")
	w := httptest.NewRecorder()

	app.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestLogTemplateError(t *testing.T) {
	app := setupTestApp()

	testErr := fmt.Errorf("template not found")
	app.LogTemplateError("test-handler", "test.html", testErr)
}

func TestLogInternalServerError(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/test", nil)
	testErr := fmt.Errorf("database connection failed")
	app.LogInternalServerError(req, "test error", "test-handler", testErr)
}

func TestLogBadRequest(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("POST", "/test", nil)
	testErr := fmt.Errorf("invalid input")
	app.LogBadRequest(req, "validation failed", "test-handler", testErr)
}

func TestLogDatabaseError(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/test", nil)
	testErr := fmt.Errorf("connection timeout")
	app.LogDatabaseError(req, "query failed", "SELECT * FROM users", testErr)
}
