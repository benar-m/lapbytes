package main

import (
	"context"
	"lapbytes/internal/api"
	"log"
	"log/slog"
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
	jsonlogger := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonlogger)

	app := &api.App{
		DB:     pool,
		Logger: logger,
	}
	// Public Routes
	mux.Handle("GET /{$}", http.HandlerFunc(app.RenderHome))
	mux.Handle("GET /register", http.HandlerFunc(app.RenderRegister))
	mux.Handle("GET /login", http.HandlerFunc(app.RenderLogin))
	mux.Handle("GET /products", http.HandlerFunc(app.RenderProducts))
	mux.Handle("GET /product/{id}", http.HandlerFunc(app.RenderProduct))

	// Auth APIs
	mux.Handle("POST /api/login", http.HandlerFunc(app.LoginUser))
	mux.Handle("POST /api/register", http.HandlerFunc(app.RegisterUser))

	// Protected User API
	mux.Handle("GET /api/product/{id}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		http.HandlerFunc(app.ListProduct),
	)))
	mux.Handle("GET /api/products/{limit}/{page}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		http.HandlerFunc(app.ListProducts),
	)))

	// Admin-only Routes
	mux.Handle("GET /api/admin/listusers/{limit}/{page}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		app.IsAdminJwtVerifierMW(
			http.HandlerFunc(app.ListUsers),
		),
	)))
	mux.Handle("GET /api/admin/listuser/{id}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		app.IsAdminJwtVerifierMW(
			http.HandlerFunc(app.ListSingleUser),
		),
	)))
	mux.Handle("POST /api/admin/deleteuser/{id}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		app.IsAdminJwtVerifierMW(
			http.HandlerFunc(app.DeleteUser),
		),
	)))
	mux.Handle("POST /api/admin/deleteproduct/{id}", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		app.IsAdminJwtVerifierMW(
			http.HandlerFunc(app.DeleteProduct),
		),
	)))
	mux.Handle("POST /api/admin/addproduct", app.ReqLoggingMW(app.GeneralJwtVerifierMW(
		app.IsAdminJwtVerifierMW(
			http.HandlerFunc(app.AddNewProduct),
		),
	)))

	// mux.HandleFunc("GET /api/admin/listusers/{limit}/{page}", app.ListUsers)
	// mux.HandleFunc("GET /api/admin/listuser/{id}", app.ListSingleUser)
	// mux.HandleFunc("POST /api/admin/deleteuser/{id}", app.DeleteUser)
	// mux.HandleFunc("POST /api/admin/deleteproduct/{id}", app.DeleteProduct)
	// mux.HandleFunc("POST /api/admin/addproduct", app.AddNewProduct)

	log.Print("Starting Server")
	err = http.ListenAndServe(":5050", mux)
	if err != nil {
		log.Fatalf("Error Starting Server %v", err)
	}

}
