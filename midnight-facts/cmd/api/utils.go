package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelope map[string]interface{}

// function to write JSON
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, header http.Header) error {
	// Encode the data to json
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// append a new line just to make it easier to read in terminal
	payload = append(payload, '\n')

	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)

	return nil
}

// function to read JSON
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// begin by setting the max bytes to limit the size of the request body
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// initialize json decoder, Calling DisallowUnknownFiels() before decoding.
	// If the JSON from the client includes any field which cannot be mapped to
	// the target destination, the decoder will return an error
	payload := json.NewDecoder(r.Body)
	payload.DisallowUnknownFields()

	// decode the request body into the target
	err := payload.Decode(dst)
	if err != nil {
		return err
	}

	// Calling Decode() again, using a pointer to an empty struct as the target.
	// If the request body contains only a single JSON value, the *will* return an io.EOF.
	// So, if we get anything else, we know the body contains additional data.
	err = payload.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only a single JSON value")
	}

	return nil
}
