package main

import (
	"log"
	"net/http"
)

// create errorResponse to be used throughout error handling
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		log.Println(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse will be used when the application encounters an error at run time.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse is used to send a 404 Not Found status code and JSON response back to the client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}
