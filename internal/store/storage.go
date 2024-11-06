package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"audio-go/internal/auth" // Make sure to import the auth package where JWT authenticator is defined
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		SignIn(context.Context, *User) (*SignInResponse, error)
		SignUp(context.Context, *User) (*SignUpResponse, error)
	}
}

// NewStorage creates a new Storage instance and initializes UserStore
func NewStorage(db *sql.DB, jwt auth.JWTAuthenticator) Storage {
	return Storage{
		Users: &UserStore{
			db:      db,
			jwtAuth: &jwt, // Pass the jwt authenticator
		},
	}
}







// package store

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"time"
// )

// var (
// 	ErrNotFound          = errors.New("resource not found")
// 	ErrConflict          = errors.New("resource already exists")
// 	QueryTimeoutDuration = time.Second * 5
// )

// type Storage struct {
// 	Users interface {
// 		SignIn(context.Context, *User)(*SignInResponse, error)
// 		SignUp(context.Context, *User) (*SignUpResponse, error)
// 	}
// }

// func NewStorage(db *sql.DB, jwt jwtAuthenticator) Storage {
// 	return Storage{
// 		Users: &UserStore{db},
// 	}
// }
