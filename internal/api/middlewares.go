package api

import "net/http"

/*
logging middleware,
jwt verification for an admin
normal jwt verifier
jwt verifier for normal browsing




*/

//inventory management (simple)
//better logging

// Is a middleware to log http requests
func ReqLoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
