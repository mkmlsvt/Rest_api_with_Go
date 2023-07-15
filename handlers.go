package main

import (
	"encoding/json"
	"fmt"
	"myProject/internal/data"
	"net/http"
	"strconv"
)

func (app *application) healtcheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     "1.0.0",
	}
	js, err := json.MarshalIndent(data, "", "\n")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
func (app *application) getCreateFootballPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		players, err := app.models.FootballPlayers.GetAll()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if err := app.writeJson(w, http.StatusOK, envelope{"players": players}, nil); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		var input struct {
			Name     string `json:"name"`
			LastName string `json:"last-name"`
			Value    int32  `json:"value"`
			Team     string `json:"team"`
		}
		err := app.readJson(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		footballPlayer := &data.FootballPlayer{
			Name:     input.Name,
			LastName: input.LastName,
			Value:    input.Value,
			Team:     input.Team}
		err = app.models.FootballPlayers.Insert(footballPlayer)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/footballPlayers/%d", footballPlayer.Id))
		err = app.writeJson(w, http.StatusOK, envelope{"player": input}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
}
func (app *application) getUpdateDeleteFootballPlayerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.getFootballPlayer(w, r)
	case http.MethodPut:
		app.updateFootballPlayer(w, r)
	case http.MethodDelete:
		app.deleteFootballPlayer(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getFootballPlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/players/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	player, err := app.models.FootballPlayers.Get(idInt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err := app.writeJson(w, http.StatusOK, envelope{"player": player}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (app *application) updateFootballPlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/players/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest)+"id hata", http.StatusBadRequest)
		return
	}

	footballPlayer, err := app.models.FootballPlayers.Get(idInt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	var input struct {
		Name     string `json:"name"`
		LastName string `json:"last-name"`
		Value    int32  `json:"value"`
		Team     string `json:"team"`
	}

	err = app.readJson(w, r, &input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if input.Team != "" {
		footballPlayer.Team = input.Team
	}
	if input.Name != "" {
		footballPlayer.Name = input.Name
	}
	if input.LastName != "" {
		footballPlayer.LastName = input.LastName
	}
	if input.Value != 0 {
		footballPlayer.Value = input.Value
	}

	err = app.models.FootballPlayers.Update(footballPlayer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.writeJson(w, http.StatusOK, envelope{"player": footballPlayer}, nil)

}
func (app *application) deleteFootballPlayer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/v1/players/"):]
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	err = app.models.FootballPlayers.Delete(idInt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	err = app.writeJson(w, http.StatusOK, envelope{"FootballPlayer": "player successfully deleted"}, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}

//getUpdateDeleteBooksHandler
