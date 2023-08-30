package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agpelkey/midnight-kittens/internal/data"
)

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






