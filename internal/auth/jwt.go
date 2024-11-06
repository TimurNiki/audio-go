package auth // Declares the package name as 'auth'

import (
	"fmt" // Imports the 'fmt' package for formatted I/O
	"time"

	"github.com/golang-jwt/jwt/v5" // Imports the JWT package for handling JSON Web Tokens
)

// Defines a struct named JWTAuthenticator with fields for secret, audience, and issuer
type JWTAuthenticator struct {
	secret string // Secret key used for signing tokens
	aud    string // Audience for the token
	iss    string // Issuer of the token
}

// Constructor function to create a new JWTAuthenticator
func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{ // Returns a pointer to a new instance of JWTAuthenticator
		secret, iss, aud, // Initializes fields with provided values (order: secret, iss, aud)
	}
}

// Method to generate a JWT token with given claims
func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Creates a new JWT token with specified signing method and claims

	// Signs the token with the secret key and converts it to a byte slice
	tokenString, err := token.SignedString([]byte(a.secret)) 
	if err != nil { // Checks for errors during signing
		return "", err // Returns an empty string and the error if signing fails
	}

	return tokenString, nil // Returns the signed token string and a nil error
}

// Method to validate a given JWT token
func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	// Parses the token and verifies its signing method using a callback function
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// Checks if the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"]) // Returns an error for unexpected signing methods
		}
		return []byte(a.secret), nil // Returns the secret key for signature verification
	},
		// Additional options for token validation
		jwt.WithExpirationRequired(), // Ensures the token has not expired
		jwt.WithAudience(a.aud), // Validates the token's audience
		jwt.WithIssuer(a.iss), // Validates the token's issuer
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), // Ensures the token uses the specified signing method
	)
}


// CreateStandardClaims creates standard JWT claims for a user
func (a *JWTAuthenticator) CreateStandardClaims(userID int64, email string, duration time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   a.iss,
		"aud":   a.aud,
		"exp":   time.Now().Add(duration).Unix(),
	}
}
