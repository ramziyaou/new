package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"time"
)

// Timer middleware with trailer
func TimerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Trailer", "execution")
		next.ServeHTTP(w, r)
		log.Printf("Request handling time: %v\n", time.Now().Sub(start))

		w.Header().Add("execution", time.Now().Sub(start).String())
	})
}

// Validate method
func HTTPMethodsCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusNotImplemented)
			fmt.Fprintln(w, "Unsupported method")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Authenticate user
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Please provide credentials")
			return
		}
		var user entity.User
		// Get username from value passed in request context upon authentication
		if err := database.Connector.Where("username = ?", username).First(&user).Error; err == nil && user.Password == password {
			ctx := context.WithValue(r.Context(), "username", username)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Invalid username or password")
		}
	})
}
