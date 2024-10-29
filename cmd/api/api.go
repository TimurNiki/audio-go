package main

import (
	"audio-go/internal/auth"
	"audio-go/internal/store"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type application struct {
	config        config
	store         store.Storage
	authenticator auth.Authenticator
}

type config struct {
	addr string
	db   dbConfig
	env  string
	auth authConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}
type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler  {
	r:=chi.NewRouter()
}


func (app *application) run(mux http.Handler) error {
	// Create a new HTTP server with specified configurations
	srv := &http.Server{
		Addr:         app.config.addr,  // Set the address to listen on from the config
		Handler:      mux,              // Assign the HTTP handler (router/mux) to the server
		WriteTimeout: time.Second * 30, // Set the maximum duration for writing the response
		ReadTimeout:  time.Second * 10, // Set the maximum duration for reading the request
		IdleTimeout:  time.Minute,      // Set the maximum idle time for connections
	}

	// Channel to signal when the server should shut down
	shutdown := make(chan error)

	// Start a goroutine to listen for termination signals
	go func() {
		// Channel to receive OS signals
		quit := make(chan os.Signal, 1)

		// Notify the quit channel on SIGINT (Ctrl+C) or SIGTERM (termination request)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a signal to be received
		s := <-quit

		// Create a context with a 5-second timeout for the shutdown process
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel() // Ensure the context is cancelled after use

		// Log that a termination signal has been caught
		app.logger.Infow("signal caught", "signal", s.String())

		// Initiate the shutdown process of the server and send the result to the shutdown channel
		shutdown <- srv.Shutdown(ctx)
	}()

	// Log that the server has started, including address and environment info
	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	// Start the server and listen for incoming requests
	err := srv.ListenAndServe()

	// Check if the error returned is due to the server being closed intentionally
	if !errors.Is(err, http.ErrServerClosed) {
		return err // If it's a different error, return it
	}

	// Wait for the shutdown signal to complete and retrieve any resulting error
	err = <-shutdown
	if err != nil {
		return err // Return any error that occurred during shutdown
	}

	// Log that the server has stopped
	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil // Indicate that the run method completed successfully
}