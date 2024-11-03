package store

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Define custom error messages
var (
	ErrUserNotFound    = errors.New("user not found")      // User does not exist
	ErrInvalidPassword = errors.New("invalid password")    // Password is incorrect
	ErrEmailTaken      = errors.New("email already taken") // Email is already registered
)

// User struct represents a user in the system
type User struct {
	ID        int64    `json:"id"`         // Unique identifier for the user
	Email     string   `json:"email"`      // User's email address
	Password  password `json:"-"`          // User's password, not exposed in JSON
	CreatedAt string   `json:"created_at"` // Timestamp when the user was created
}

// password struct encapsulates password management
type password struct {
	text *string // Plaintext password (optional)
	hash []byte  // Hashed password
}

// Set hashes the plaintext password and stores it
func (p *password) Set(text string) error {
	// Generate a bcrypt hash from the password
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err // Return error if hashing fails
	}

	p.text = &text // Store the plaintext password (for later use)
	p.hash = hash  // Store the hashed password

	return nil
}

// UserStore handles user-related database operations
type UserStore struct {
	db *sql.DB // Database connection
}

// SignIn authenticates a user with their email and password
func (us *UserStore) SignIn(ctx context.Context, user *User) error {
	var storedHash []byte
	// Query to retrieve the hashed password for the provided email
	err := us.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound // User not found
		}
		return err // Handle other errors
	}

	// Ensure the password text is set
	if user.Password.text == nil {
		return ErrInvalidPassword // Password is not set
	}

	// Compare the stored hash with the provided plaintext password
	if err := bcrypt.CompareHashAndPassword(storedHash, []byte(*user.Password.text)); err != nil {
		return ErrInvalidPassword // Password mismatch
	}

	return nil // Successful sign-in
}

// SignUp registers a new user
func (us *UserStore) SignUp(ctx context.Context, user *User) error {
	// Check if the user already exists by querying their email
	var existingID int64
	err := us.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingID)
	if err == nil {
		return ErrEmailTaken // Email is already taken
	} else if err != sql.ErrNoRows {
		return err // Handle unexpected database error
	}
	// Ensure the password is set before signing up
	if user.Password.text == nil {
		return errors.New("password must be set before signing up") // Password not provided
	}

	// Hash the password before saving to the database
	if err := user.Password.Set(*user.Password.text); err != nil {
		return err // Return error if hashing fails
	}

	// Insert the new user into the database
	_, err = us.db.ExecContext(ctx, "INSERT INTO users (email, password, created_at) VALUES (?, ?, NOW())", user.Email, user.Password.hash)
	return err // Return any error from the insertion
}
