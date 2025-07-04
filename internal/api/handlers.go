package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lapbytes/internal/model"
	"lapbytes/internal/store/queries"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB     *pgxpool.Pool
	Logger *slog.Logger
}

// RenderHome serves the homepage template
func (a *App) RenderHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		a.LogInternalServerError(r, "template parsing", "renderhome", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		a.LogTemplateError("renderhome", "index.gohtml", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

// RenderRegister serves the user registration page
func (a *App) RenderRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.gohtml")
	if err != nil {
		a.LogInternalServerError(r, "template parsing", "renderregister", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		a.LogTemplateError("renderregister", "signup.gohtml", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// RenderLogin serves the user login page
func (a *App) RenderLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.gohtml")
	if err != nil {
		a.LogInternalServerError(r, "template parsing", "renderlogin", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html ; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		a.LogTemplateError("renderlogin", "login.html", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// RenderProducts serves the products listing page
func (a *App) RenderProducts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		a.LogInternalServerError(r, "template parsing", "renderproducts", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		a.LogTemplateError("renderproducts", "index.gohtml", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
}

// RenderProduct serves the individual product details page
func (a *App) RenderProduct(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/product-details.gohtml")
	if err != nil {
		a.LogInternalServerError(r, "template parsing", "renderproduct", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		a.LogTemplateError("renderproduct", "product-details.gohtml", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// LoginUser authenticates a user and issues JWT token
func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		a.LogBadRequest(r, "invalid content-type", "loginuser", fmt.Errorf("non json request"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bad request, accepts JSON only",
		})
		return

	}
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var userRequest loginRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		a.LogBadRequest(r, "possibly missing values", "loginuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "bad request, All fields are required",
		})
		return
	}

	passwordhash, err := queries.GetUserHash(a.DB, userRequest.Email)
	if err != nil {
		a.Logger.Error("invalid credentials",
			"handler", "loginuser",
			"path", r.URL.Path,
			"method", r.Method,
			"status", 401,
			"error", fmt.Errorf("user does not exist"),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Credentials. Please Check again",
		})

		return
	}

	loggedIn := verifyPasswordHash(userRequest.Password, passwordhash)
	if !loggedIn {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			a.Logger.Error("ip parsing err",
				"remoteAddress", r.RemoteAddr,
			)
		}
		a.Logger.Error("invalid credentials",
			"handler", "loginuser",
			"path", r.URL.Path,
			"method", r.Method,
			"status", 401,
			"ip", host,
			"error", fmt.Errorf("invalid password"),
		)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Credentials. Please Check again",
		})

		return
	}

	InitKeys()
	accessToken, err := IssueKeys()
	if err != nil {
		a.LogInternalServerError(r, "jwt token issuing error", "loginuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An Error Occured During Login",
		})
		return
	}

	response := model.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	//Set cookies + A refresh token
	refreshToken, err := generateRandomString()
	if err != nil {
		a.LogInternalServerError(r, "cookie issuance error", "loginuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An Error Occured During Login",
		})
		return

	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(time.Hour * 24 * 3),
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error Writing Response Body %+v", err)
	}
}

// RegisterUser creates a new user account
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//later sanitize, middleware maybe
	//refactor to more direct or keep it as defensive as it is
	type regRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if r.Header.Get("Content-Type") != "application/json" {
		a.LogBadRequest(r, "invalid content-type", "registeruser", fmt.Errorf("invalid content-type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid data, please check again",
		})
		return
	}
	var userRequest regRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		a.LogBadRequest(r, "probably missing values", "registeruser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid data, please check again",
		})
		return

	}

	user := &model.User{}
	password_hash, err := hashPassword(userRequest.Password)
	if err != nil {
		a.LogInternalServerError(r, "password hashing error", "registeruser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error occured during registration",
		})
		return

	}

	current_time := time.Now()

	user.Username = userRequest.Username
	user.Email = userRequest.Email
	user.Password_hash = password_hash
	user.Is_admin = false
	user.Access_level = 4
	user.Created_at = current_time
	user.Updated_at = current_time

	userId, err := queries.InsertUser(a.DB, *user)
	if err != nil {
		a.LogDatabaseError(r, "insert query error", "insertuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Could not Create Account. Email Might already exist",
		})
		return
	}
	a.Logger.Info("successful user creation",
		"time", time.Now(),
		"userid", userId,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "account created successfully",
	})

}

// ListProducts returns paginated list of all laptops
func (a *App) ListProducts(w http.ResponseWriter, r *http.Request) {
	//lojik to get the Products
	limit := r.PathValue("limit")
	page := r.PathValue("page")
	lim, err := strconv.Atoi(limit)
	if err != nil || lim < 1 {
		a.LogBadRequest(r, "invalid limit", "listproducts", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Request",
		})
		return
	}
	pag, err := strconv.Atoi(page)
	if err != nil || pag < 1 {
		a.LogBadRequest(r, "invalid page", "listproducts", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Request",
		})
		return
	}
	offset := (pag - 1) * lim

	products, err := queries.QueryLaptops(a.DB, lim, offset)
	if err != nil {
		a.LogDatabaseError(r, "query laptops error", "querylaptops", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return

	}
	a.Logger.Info("successful products query",
		"time", time.Now(),
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
		"message":  "request successful",
	})

}

// ListProduct returns details of a single laptop by ID
func (a *App) ListProduct(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")
	id, err := strconv.Atoi(productId)
	if err != nil || id < 1 {
		a.LogBadRequest(r, "invalid product request id", "listproduct", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "item not found",
		})
		return
	}
	product, err := queries.QueryLaptop(a.DB, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			a.Logger.Error("item not found",
				"status", 404,
				"id", id,
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "item not found",
			})
			return
		}
		a.LogDatabaseError(r, "product query error", "querylaptop", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product) //Maybe encode to memory first later to avoid sending malformed json
	if err != nil {
		a.Logger.Error("encoding error",
			"handler", "querylaptop",
			"error", err,
		)
		return

	}
	a.Logger.Info("product query was successul",
		"time", time.Now(),
	)

}

// DeleteUser removes a user from the database (admin only)
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("id")
	id, err := strconv.Atoi(user_id)
	if err != nil || id < 1 {
		a.LogBadRequest(r, "invalid user id", "deleteuser", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user id",
		})
		return
	}
	err = queries.DeleteUser(a.DB, id)
	if err != nil {
		a.LogDatabaseError(r, "delete user query error", "deleteuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return

	}
	a.Logger.Info("successful user deletion",
		"id", id,
		"time", time.Now(),
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user successfully deleted",
		"userid":  id,
	})
}

// ListUsers returns paginated list of all users (admin only)
func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.PathValue("limit")
	page := r.PathValue("page")

	lim, err := strconv.Atoi(limit)
	if err != nil || lim < 1 {
		a.LogBadRequest(r, "invalid limit", "listusers", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid limit",
		})
		return
	}
	pag, err := strconv.Atoi(page)
	if err != nil || pag < 1 {
		a.LogBadRequest(r, "invalid page", "listusers", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid page number",
		})
		return
	}
	offset := (pag - 1) * lim

	users, err := queries.GetAllUsers(a.DB, lim, offset)
	if err != nil {
		a.LogDatabaseError(r, "get all users query error", "listusers", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	a.Logger.Info("successful users listing",
		"time", time.Now(),
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "request successful",
		"users":   users,
	})

}

// ListSingleUser returns details of a specific user by ID (admin only)
func (a *App) ListSingleUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("id")
	id, err := strconv.Atoi(user_id)
	if err != nil || id < 1 {
		a.LogBadRequest(r, "invalid user id", "listsingleuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user id",
		})
		return
	}
	user, err := queries.GetUser(a.DB, id)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with id %d not found", id) {
			a.LogDatabaseError(r, "user not found", "listsingleuser", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user id does not exist",
			})
			return

		}

		a.LogDatabaseError(r, "get user query error", "listsingleuser", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	a.Logger.Info("successful user listing",
		"time", time.Now(),
	)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "request successful",
		"user":    user,
	})

}

// AddNewProduct creates a new laptop in the database (admin only)
func (a *App) AddNewProduct(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		a.LogBadRequest(r, "invalid content-type", "addnewproduct", fmt.Errorf("non json request"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid data type",
		})
		return
	}

	var product model.Laptop
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		a.LogBadRequest(r, "invalid json", "addnewproduct", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON",
		})
		return
	}
	productId, err := queries.InsertLaptop(a.DB, product)
	if err != nil {
		a.LogDatabaseError(r, "insert laptop query error", "insertlaptop", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "could not add laptop",
		})
		return
	}
	a.Logger.Info("laptop successfully added",
		"time", time.Now(),
		"id", productId,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product added successfully",
		"id":      productId,
	})

}

// DeleteProduct removes a laptop from the database (admin only)
func (a *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil || productId < 1 {
		a.LogBadRequest(r, "invalid id", "deleteproduct", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}
	err = queries.DeleteLaptop(a.DB, productId)
	if err != nil {
		a.LogDatabaseError(r, "delete laptop query error", "deleteproduct", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	a.Logger.Info("product successfully deleted",
		"time", time.Now(),
		"id", productId,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product deleted successfully",
		"id":      productId,
	})

}
