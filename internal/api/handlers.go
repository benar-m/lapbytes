package api

import (
	"encoding/json"
	"fmt"
	"lapbytes/internal/model"
	"lapbytes/internal/store/queries"
	"log"
	"net/http"
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
	fmt.Fprint(w, "Hello Home Page")
}

func (a *App) RenderRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Register Page")
}
func (a *App) RenderLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Login Page")
}
func (a *App) RenderProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Products Page")
}

func (a *App) RenderProduct(w http.ResponseWriter, r *http.Request) {
	prod_id := r.PathValue("id")
	fmt.Fprintf(w, "Hello from Product id %+v", prod_id)
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
	passwordhash, err := queries.GetUserHash(a.DB, email)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	loggedIn := verifyPasswordHash(password, passwordhash)
	if !loggedIn {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	InitKeys()
	accessToken, err := IssueKeys()
	if err != nil {
		http.Error(w, "Error Occured During Login", http.StatusInternalServerError)
		return
	}
	response := model.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	//Set cookies + A refresh token
	refreshToken, err := generateRandomString()
	if err != nil {
		http.Error(w, "Error Occured During Login", http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Error Writing Response Body %+v", err)
	}
}

// Called at /api/..., this function creates a new user in the database
func (a *App) RegisterUser(w http.ResponseWriter, r *http.Request) {
	//later sanitize, middleware maybe
	//refactor to more direct or keep it as defesive as it is
	user := &model.User{}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password_hash, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Error Occured During Registration", http.StatusInternalServerError)
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
		http.Error(w, "Could not Create Account", http.StatusInternalServerError)
		return
	}
	log.Print(userId)

}
