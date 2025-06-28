package main

import (
	"lapbytes/internal/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("static/", fileServer))
	mux.HandleFunc("GET /{$}", api.Home)
	mux.HandleFunc("GET /register", api.Register)
	mux.HandleFunc("GET /login", api.Login)
	mux.HandleFunc("GET /products", api.Products)
	mux.HandleFunc("GET /product", api.Product)

	mux.HandleFunc("POST api/login", api.ApiLogin)

	log.Print("Starting Server")
	err := http.ListenAndServe(":5050", mux)
	if err != nil {
		log.Fatalf("Error Starting Server %v", err)
	}

}
