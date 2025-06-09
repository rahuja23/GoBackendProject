package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		_ = writeJSONError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
	}
}

func (app *application) readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ready!!"))
}

func (app *application) buildInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Build Info"))
}

func (app *application) liveHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Live"))
}
