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
	mux.Handle("/static/", http.StripPrefix("static/", fileServer))
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
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("GET /register", app.Register)
	mux.HandleFunc("GET /login", app.Login)
	mux.HandleFunc("GET /products", app.Products)
	mux.HandleFunc("GET /product", app.Product)

	mux.HandleFunc("POST api/login", app.ApiLogin)

	log.Print("Starting Server")
	err = http.ListenAndServe(":5050", mux)
	if err != nil {
		log.Fatalf("Error Starting Server %v", err)
	}

}
