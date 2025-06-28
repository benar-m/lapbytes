package api

import (
	"encoding/json"
	"fmt"
	"lapbytes/internal/model"
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
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Home Page")
}

func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Register Page")
}
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Login Page")
}
func (a *App) Products(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Products Page")
}

func (a *App) Product(w http.ResponseWriter, r *http.Request) {
	prod_id := r.PathValue("id")
	fmt.Fprintf(w, "Hello from Product id %+v", prod_id)
}

// func ApiLogin(w http.ResponseWriter, r *http.Request) {
// 	// username:=r.FormValue("username")
// 	// password:=r.FormValue("password")
// 	//Verify Login Then Issue KEys
// }

func (a *App) ApiLogin(w http.ResponseWriter, r *http.Request) {

	// username:=r.FormValue("email")
	// password:=r.FormValue("password")
	InitKeys()
	accessToken, err := IssueKeys()
	if err != nil {
		http.Error(w, "Error Occured During Login", http.StatusInternalServerError)
	}
	response := model.LoginResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}

	//Set cookies + A refresh token
	refreshToken, err := generateRandomString()
	if err != nil {
		http.Error(w, "Error Occured During Login", http.StatusInternalServerError)

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
