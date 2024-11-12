package main

import (
	"audio-go/internal/store" // Adjust this import path according to your project structure
	"context"
	"net/http"
)

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
func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Use NewPassword to initialize the password field
	user := &store.User{
		Email:    req.Email,
		Password: *store.NewPassword(req.Password), // Initialize password with NewPassword constructor
	}

	ctx := context.Background()
	// Call SignIn method from the store
	resp, err := app.store.Users.SignIn(ctx, user)
	if err != nil {
		// Handle errors based on the type
		switch err {
		case store.ErrUserNotFound, store.ErrInvalidPassword:
			app.badRequestResponse(w, r, err) // Bad request for invalid user or password
		default:
			app.internalServerError(w, r, err) // Internal error for other cases
		}
		return
	}

	// Send the successful response
	writeJSON(w, http.StatusOK, resp) // Sending JWT token + user info back
}

// SignUp handles user sign-up
func (app *application) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := readJSON(w, r, &req); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Use NewPassword to initialize the password field
	user := &store.User{
		Email:    req.Email,
		Password: *store.NewPassword(req.Password), // Initialize password with NewPassword constructor
	}
	ctx := context.Background()
	// Call SignUp method from the store
	resp, err := app.store.Users.SignUp(ctx, user)
	if err != nil {
		// Handle errors based on the type
		switch err {
		case store.ErrEmailTaken, store.ErrPasswordNotSet:
			app.badRequestResponse(w, r, err) // Bad request for duplicate email or empty password
		default:
			app.internalServerError(w, r, err) // Internal error for other cases
		}
		return
	}

	// Send the successful response
	writeJSON(w, http.StatusCreated, resp) // Created status for new user registration
}
