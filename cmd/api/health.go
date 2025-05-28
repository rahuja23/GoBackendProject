package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
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
