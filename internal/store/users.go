package store

import (
	"audio-go/internal/auth"
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrEmailTaken       = errors.New("email already taken")
	ErrPasswordNotSet   = errors.New("password must be set before signing up")
)

// User represents a user in the system
type User struct {
	ID        int64    `json:"id"`
	Email     string   `json:"email"`
	Password  password `json:"-"` // Unexported password field (we don't expose it in the response)
	CreatedAt string   `json:"created_at"`
}

// password manages password hashing and verification
type password struct {
	text *string // Plaintext password (for comparison only)
	hash []byte  // The hashed version of the password
}

// NewPassword creates and returns a new password object initialized with the given plaintext password
func NewPassword(text string) *password {
	return &password{
		text: &text,
	}
}

// Set hashes the plaintext password and stores the hash
func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store the plaintext password (optional) and the hashed version
	p.text = &text
	p.hash = hash

	return nil
}

// Compare compares the given plaintext password with the stored hash
func (p *password) Compare(plainText string) error {
	// If text is nil, we can't compare it (no password has been set)
	if p.text == nil {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
}

// UserStore handles user-related database operations
type UserStore struct {
	db      *sql.DB
	jwtAuth *auth.JWTAuthenticator
}

// NewUserStore creates a new UserStore
func NewUserStore(db *sql.DB, jwtAuth *auth.JWTAuthenticator) *UserStore {
	return &UserStore{
		db:      db,
		jwtAuth: jwtAuth,
	}
}

// SignInResponse represents the response returned on successful sign-in
type SignInResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// SignIn authenticates a user and returns a JWT token if successful
func (us *UserStore) SignIn(ctx context.Context, user *User) (*SignInResponse, error) {
	var storedHash []byte
	// Retrieve stored password hash from the database
	err := us.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ?", user.Email).Scan(&user.ID, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Ensure the provided password is not empty
	if user.Password.text == nil {
		return nil, ErrInvalidPassword
	}

	// Compare the stored password hash with the input password
	if err := bcrypt.CompareHashAndPassword(storedHash, []byte(*user.Password.text)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Create JWT claims
	claims := us.jwtAuth.CreateStandardClaims(user.ID, user.Email, time.Hour*24) // Token expires in 24 hours

	// Generate JWT token
	token, err := us.jwtAuth.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// Return token and user details
	return &SignInResponse{
		Token: token,
		User:  user,
	}, nil
}

// SignUpResponse represents the response returned on successful sign-up
type SignUpResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// SignUp registers a new user and returns a JWT token if successful
func (us *UserStore) SignUp(ctx context.Context, user *User) (*SignUpResponse, error) {
	// Check if the user already exists
	var existingID int64
	err := us.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingID)
	if err == nil {
		// Email already taken
		return nil, ErrEmailTaken
	} else if err != sql.ErrNoRows {
		// Some other database error
		return nil, err
	}

	// Ensure the password is set before signing up
	if user.Password.text == nil {
		return nil, ErrPasswordNotSet
	}

	// Hash the password before storing it
	if err := user.Password.Set(*user.Password.text); err != nil {
		return nil, err
	}

	// Insert the new user into the database
	result, err := us.db.ExecContext(ctx, "INSERT INTO users (email, password, created_at) VALUES (?, ?, NOW())", user.Email, user.Password.hash)
	if err != nil {
		return nil, err
	}

	// Get the new user's ID
	user.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Define claims for the JWT token
	claims := us.jwtAuth.CreateStandardClaims(user.ID, user.Email, time.Hour*24) // Token expires in 24 hours

	// Generate JWT token
	token, err := us.jwtAuth.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// Return token and user details
	return &SignUpResponse{
		Token: token,
		User:  user,
	}, nil
}
