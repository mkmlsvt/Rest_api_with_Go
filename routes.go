package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healtcheck)
	mux.HandleFunc("/v1/players", app.getCreateFootballPlayerHandler)
	mux.HandleFunc("/v1/players/", app.getUpdateDeleteFootballPlayerHandler)
	return mux
}
