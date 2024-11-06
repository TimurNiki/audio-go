// handlers/middleware.go
package main

import (
	"audio-go/internal/auth"
	"fmt"
	"net/http"
)

// AuthMiddleware validates JWT tokens for protected routes
func AuthMiddleware(jwtAuth *auth.JWTAuthenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Validate the token
		parsedToken, err := jwtAuth.ValidateToken(token)
		if err != nil || !parsedToken.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Optionally, set user info in context here

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
