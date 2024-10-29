package main

import (
	"encoding/json"        // Importing the encoding/json package for JSON encoding/decoding
	"net/http"            // Importing the net/http package for HTTP server and client implementations
	"github.com/go-playground/validator/v10" // Importing the validator package for struct validation
)

// Validate is a global variable to hold the validator instance
var Validate *validator.Validate

// init function is automatically called when the package is initialized
func init() {
	// Create a new validator instance with required struct validation enabled
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

