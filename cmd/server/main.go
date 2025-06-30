package main

import (
	"context"
	"lapbytes/internal/api"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	mux := http.NewServeMux()
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	dsn := os.Getenv("PG_DATABASE_URL")
	if dsn == "" {
		log.Fatal("No database")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to Connect to the database: %+v", err)
	}
	defer pool.Close()

	app := &api.App{
		DB: pool,
	}
	mux.HandleFunc("GET /{$}", app.RenderHome)
	mux.HandleFunc("GET /register", app.RenderRegister)
	mux.HandleFunc("GET /login", app.RenderLogin)
	mux.HandleFunc("GET /products", app.RenderProducts)
	mux.HandleFunc("GET /product", app.RenderProduct)

	mux.HandleFunc("POST /api/login", app.LoginUser)
	mux.HandleFunc("POST /api/register", app.RegisterUser)
	mux.HandleFunc("POST /api/product/{id}", app.ListProduct)
	loggedHandler := api.RouteLogger(mux)

	log.Print("Starting Server")
	err = http.ListenAndServe(":5050", loggedHandler)
	if err != nil {
		log.Fatalf("Error Starting Server %v", err)
	}

}
