package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agpelkey/midnight-kittens/internal/data"
)

/*
func (app *application) handleStoreCatFact(w http.ResponseWriter, r *http.Request) {

    payload, err := app.handleGetCatFact()
    if err != nil {
        return
    }

    fact := strings.NewReader(payload.Fact)


    req, err := http.NewRequest("POST", "https://localhost/api/getFacts", fact)
    if err != nil {
        return
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return
    }

    defer resp.Body.Close()

    err = app.Models.Facts.SendFactToDB(payload)

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("api/getFacts/%d", payload.Id))

    err = app.writeJSON(w, http.StatusCreated, envelope{"fact": fact}, headers)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }
}
*/

func (app *application) handleGetCatFact() (*data.CatFact, error) {
	req, err := http.NewRequest("GET", "https://catfact.ninja/fact", nil)
	if err != nil {
		return nil, err
	}

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close()

    var input struct {
        Fact string `json:"fact"`
        Length int `json:"length"`
    }

    err = json.NewDecoder(resp.Body).Decode(&input)
    if err != nil {
        return nil, err
    }

    result := &data.CatFact{
        Fact: input.Fact,
        Length: input.Length,
    }


    fmt.Println(result)
    return result, nil
    
}






