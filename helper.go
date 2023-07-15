package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelope map[string]any

func (app *application) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\n")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJson(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_579
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single json object")
	}
	return nil
}
