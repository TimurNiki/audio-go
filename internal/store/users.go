package store

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailTaken      = errors.New("email already taken")
)

type User struct {
	ID        int64    `json:"id"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (us *UserStore) SignIn(ctx context.Context, user *User) error {
	var storedHash []byte
	err := us.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	// Dereference the password text
	if user.Password.text == nil {
		return ErrInvalidPassword
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword(storedHash, []byte(*user.Password.text)); err != nil {
		return ErrInvalidPassword
	}

	return nil
}

func (us *UserStore) SignUp(ctx context.Context, user *User) error {
	// Check if the user already exists
	var existingID int64
	err := us.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingID)
	if err == nil {
		return ErrEmailTaken
	} else if err != sql.ErrNoRows {
		return err // Handle unexpected database error
	}
	// Hash the password before saving
	if user.Password.text == nil {
		return errors.New("password must be set before signing up")
	}

	if err := user.Password.Set(*user.Password.text); err != nil {
		return err
	}

	// Insert the new user into the database
	_, err = us.db.ExecContext(ctx, "INSERT INTO users (email, password, created_at) VALUES (?, ?, NOW())", user.Email, user.Password.hash)
	return err
}
