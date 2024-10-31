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

// writeJSON writes a JSON response with a given status code and data
func writeJSON(w http.ResponseWriter, status int, data any) error {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Write the HTTP status code
	w.WriteHeader(status)
	// Encode the data into JSON and write it to the response
	return json.NewEncoder(w).Encode(data)
}

// readJSON reads and decodes JSON data from the request body into the provided data structure
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // Set maximum bytes to read from the request body (1 MB)
	// Limit the request body size to prevent large payloads
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a new JSON decoder for the request body
	decoder := json.NewDecoder(r.Body)

	// Disallow unknown fields in the incoming JSON to enforce strict validation
	decoder.DisallowUnknownFields()

	// Decode the JSON data into the provided data structure
	return decoder.Decode(data)
}

// writeJSONError writes a JSON response with an error message and status code
func writeJSONError(w http.ResponseWriter, status int, message string) error {
	// Define an envelope structure to hold the error message
	type envelope struct {
		Error string `json:"error"` // JSON field name for the error message
	}
	// Call writeJSON to send the error response
	return writeJSON(w, status, &envelope{Error: message})
}

// jsonResponse is a method on the application struct that writes a successful JSON response
func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	// Define an envelope structure to hold the data
	type envelope struct {
		Data any `json:"data"` // JSON field name for the data
	}
	// Call writeJSON to send the data response
	return writeJSON(w, status, &envelope{Data: data})
}