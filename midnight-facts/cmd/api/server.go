package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serve() will be the method to start the http.server.
// The purpose of pulling this out of func main() is to:
//  1. keep func main() a little less cluttered.
//  2. Allow us to "gracefully shutdown the server".
func (app *application) serve() error {
	// set srv config

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// create a shudown channel. This channel is to be used to
	// catch any errors returned by the Shutdown() function.
	shutdownError := make(chan error)

	// start a background go routine
	go func() {
		// create a quit channel
		quit := make(chan os.Signal, 1)

		// Use signal.Notify to listen for incoming SIGINT and SIGTERM signals
		// and relay them to the quit channel.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Read the signal from the quit channel
		s := <-quit

		log.Println("shutting down server", map[string]string{
			"signal": s.String(),
		})

		// context bb
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Now, we call Shutdown() on our server, using the ctx we just made.
		// Shutdown() returns nil if graceful shutdown was successful, or an error if not.
		shutdownError <- srv.Shutdown(ctx)
	}()

	log.Println("starting server on port", app.config.port)

	// the important note here is that Shutdown() will cause ListenAndServer()
	// to return an http.ErrServerClosed error. So, if we can look only for this error,
	// and it tells us that the server *did* gracefully shutdown.
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to get a return value from Shutdown().
	// If this value is an error, we know that there was a problem
	// with the gracefule shutdown.
	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Println("stopped server")

	return nil
}
