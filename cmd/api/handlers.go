package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/agpelkey/midnight-kittens/internal/data"
)

func (app *application) handleStoreCatFact(w http.ResponseWriter, r *http.Request) {

}

func (app *application) handleGetCatFact() ([]*data.CatFact, error) {
	resp, err := http.Get("http://cat-fact.herokuapp.com/facts")
	if err != nil {
		return nil, err
	}

	var cf []*data.CatFact

	if err := json.NewDecoder(resp.Body).Decode(&cf); err != nil {
		return nil, err
	}

	fmt.Printf("%+v", cf)

	return cf, nil

}
