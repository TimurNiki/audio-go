package main

import (
	"audio-go/internal/auth"
	"context"
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
)

// Key for storing user data in context
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware validates JWT tokens for protected routes
// Bu fonksiyonun alıcı olarak *application türünü kullanıyoruz
func (app *application) AuthMiddleware(jwtAuth *auth.JWTAuthenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w, r, fmt.Errorf("missing Authorization header"))
			return
		}

		// Expected format: "Bearer <token>"
		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid Authorization header format"))
			return
		}

		// Validate the token
		parsedToken, err := jwtAuth.ValidateToken(token)
		if err != nil || !parsedToken.Valid {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid token"))
			return
		}

		// Extract user information from token (sub, email, etc.)
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid token claims"))
			return
		}

		// Optionally, you could fetch more user details from the database if needed
		userID, ok := claims["sub"].(float64)
		if !ok {
			app.unauthorizedResponse(w, r, fmt.Errorf("invalid token payload"))
			return
		}

		// Add the user information to the request context
		ctx := context.WithValue(r.Context(), userContextKey, userID)

		// Proceed to the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
