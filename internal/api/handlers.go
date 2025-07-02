package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"lapbytes/internal/model"
	"lapbytes/internal/store/queries"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB *pgxpool.Pool
}

/*- `GET /api/products` — List all laptops
- `GET /api/products/{id}` — Get product details
- `GET /api/products/search?q=` — Search or filter products*/

/*
## Render Endpoints
- `GET /` — Homepage
- `GET /login` — Login page
- `GET /register` — Registration page
- `GET /products` — Product listing page
- `GET /products/{id}` — Product details page
- `GET /cart` — View cart
- `GET /checkout` — Checkout page
- `GET /orders` — User’s order history
- `GET /admin` — Admin dashboard
- `GET /admin/products` — Admin laptop listing
- `GET /admin/orders` — Admin order view
- `GET /admin/users` — Admin user management

*/
// Core Rendering
func (a *App) RenderHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Printf("Error Encountered while Parsing index template : %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error Executing Index Templates")
	}
	//
}

func (a *App) RenderRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/signup.gohtml")
	if err != nil {
		log.Printf("Error while parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error While Executing Template %v", err)
		return
	}
}
func (a *App) RenderLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.gohtml")
	if err != nil {
		log.Printf("Error Encountered while parsing login template : %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html ; charset=utf-8")
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error Executing Login Template %v", err)
		return
	}
}

func (a *App) RenderProducts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.gohtml")
	if err != nil {
		log.Printf("Error Encountered while Parsing index template : %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error Executing Index Templates")
	}
}

func (a *App) RenderProduct(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/product-details.gohtml")
	if err != nil {
		log.Printf("Error while Parsing File %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error Executing Pd template: %v", err)
	}
}

// func ApiLogin(w http.ResponseWriter, r *http.Request) {
// 	// username:=r.FormValue("username")
// 	// password:=r.FormValue("password")
// 	//Verify Login Then Issue KEys
// }

// called at /api/..., it logs in a user by verifying credentials and issuing jwts
func (a *App) LoginUser(w http.ResponseWriter, r *http.Request) {
	//Logging later
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "All Fields are required",
		})
	}
	passwordhash, err := queries.GetUserHash(a.DB, email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Credentials. Please Check again",
		})
		return
	}

	loggedIn := verifyPasswordHash(password, passwordhash)
	if !loggedIn {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Credentials. Please Check again",
		})
		// http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	InitKeys()
	accessToken, err := IssueKeys()
	if err != nil {
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
		Path:     "GET /",
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

// Called at /api/..., this function creates a new user in the database
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//later sanitize, middleware maybe
	//refactor to more direct or keep it as defensive as it is
	user := &model.User{}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	if username == "" || email == "" || password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "All fields are required",
		})
		return
	}
	password_hash, err := hashPassword(password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error occured during registration",
		})
		return

	}

	current_time := time.Now()

	user.Username = username
	user.Email = email
	user.Password_hash = password_hash
	user.Is_admin = false
	user.Access_level = 4
	user.Created_at = current_time
	user.Updated_at = current_time

	userId, err := queries.InsertUser(a.DB, *user)
	if err != nil {
		log.Printf("Database Error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Could not Create Account. Email Might already exist",
		})
		return
	}
	log.Printf("created user: %v", userId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "account created successfully",
	})

}

/*
##  Products
- `GET /api/products` — List all laptops
- `GET /api/products/{id}` — Get product details
- `GET /api/products/search?q=` — Search or filter products
*/

func (a *App) ListProducts(w http.ResponseWriter, r *http.Request) {
	//lojik to get the Products
	limit := r.PathValue("limit")
	page := r.PathValue("page")
	lim, err := strconv.Atoi(limit)
	if err != nil || lim < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid Request",
		})
		return
	}
	pag, err := strconv.Atoi(page)
	if err != nil || pag < 1 {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"products": products,
	})

}

func (a *App) ListProduct(w http.ResponseWriter, r *http.Request) {
	productId := r.PathValue("id")
	id, err := strconv.Atoi(productId)
	if err != nil || id < 1 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	product, err := queries.QueryLaptop(a.DB, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			http.Error(w, "Item Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Printf("Error occured: %+v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product) //Maybe encode to memory first later to avoid sending malformed json
	if err != nil {
		//log better
		return

	}
	fmt.Println("Data Sent Succesfully")

}

//Admin endpoints

// Handle user not found well later
func (a *App) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("id")
	id, err := strconv.Atoi(user_id)
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid user id",
		})
		return
	}
	err = queries.DeleteUser(a.DB, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "user succesfully deleted",
		"userid":  id,
	})
}

func (a *App) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := r.PathValue("limit")
	page := r.PathValue("page")

	lim, err := strconv.Atoi(limit)
	if err != nil || lim < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid limit",
		})
		return
	}
	pag, err := strconv.Atoi(page)
	if err != nil || pag < 1 {

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "request successful",
		"users":   users,
	})

}

func (a *App) ListSingleUser(w http.ResponseWriter, r *http.Request) {
	user_id := r.PathValue("id")
	id, err := strconv.Atoi(user_id)
	if err != nil || id < 1 {
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user id does not exist",
			})
			return

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "request successful",
		"user":    user,
	})

}

func (a *App) AddNewProduct(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid JSON",
		})
		return
	}
	productId, err := queries.InsertLaptop(a.DB, product)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "could not add laptop",
		})
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product added successfuly",
		"id":      productId,
	})

}

func (a *App) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	productId, err := strconv.Atoi(id)
	if err != nil || productId < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}
	err = queries.DeleteLaptop(a.DB, productId)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "product deleted successfully",
		"id":      productId,
	})

}
