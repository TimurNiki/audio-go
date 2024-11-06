// main/auth.go
package main

import (
	"encoding/json"
	"net/http"
	"context"
	"audio-go/internal/store" // Adjust this import path according to your project structure
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	UserStore *store.UserStore
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userStore *store.UserStore) *AuthHandler {
	return &AuthHandler{
		UserStore: userStore,
	}
}

// SignInRequest represents the expected payload for sign-in
type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpRequest represents the expected payload for sign-up
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignIn handles user sign-in
func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use NewPassword to initialize the password field
	user := &store.User{
		Email: req.Email,
		Password: *store.NewPassword(req.Password), // Initialize password with NewPassword constructor
	}

	resp, err := h.UserStore.SignIn(context.Background(), user)
	if err != nil {
		switch err {
		case store.ErrUserNotFound, store.ErrInvalidPassword:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SignUp handles user sign-up
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use NewPassword to initialize the password field
	user := &store.User{
		Email: req.Email,
		Password: *store.NewPassword(req.Password), // Initialize password with NewPassword constructor
	}

	resp, err := h.UserStore.SignUp(context.Background(), user)
	if err != nil {
		switch err {
		case store.ErrEmailTaken, store.ErrPasswordNotSet:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
