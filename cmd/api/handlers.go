package main

import (
	"io"
	"net/http"

	"github.com/agpelkey/midnight-kittens/internal/data"
)

/*
func (app *application) handleStoreCatFact(w http.ResponseWriter, r *http.Request) {

    var input struct {
        Fact string `json:"fact"`
    }

    err := app.readJSON(w, r, &input)
    if err != nil {
    fmt.Println()
        return
    }

    fact := &data.CatFact{
        Fact: input.Fact,
    }

    err = app.Models.Facts.SendFactToDB(fact)

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("api/facts/%d", fact.Id))

    err = app.writeJSON(w, http.StatusCreated, envelope{"fact": fact}, headers)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }

}
*/

func (app *application) handleGetCatFact() error {
	req, err := http.NewRequest("GET", "https://catfact.ninja/fact", nil)
	if err != nil {
		return err
	}

    client := &http.Client{}
    resp, err := client.Do(req)
    
    defer resp.Body.Close()

    byteReponse, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // var cf *data.CatFact

    /*
    if err := json.NewDecoder(resp.Body).Decode(&cf); err != nil {
        return err
    }

    fmt.Printf("%+v", cf)
    */

    var input struct {
        Fact string `json:"fact"`
        Length string `json:"length"`
    }

    res := &data.CatFact{
        Fact: ,
        Length: ,
    }

	//fmt.Printf("%T %s\n", cf)
    err = app.Models.Facts.SendFactToDB(byteReponse)
    if err != nil {
        app.serverErrorResponse(nil, nil, err)
    }

    /*
    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("api/facts/%d", cf))

    err = app.writeJSON(nil, http.StatusCreated, envelope{"fact": cf}, headers)
    if err != nil {
        return err
    }
    */

    return nil
}







