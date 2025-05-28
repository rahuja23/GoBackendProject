package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Ready!!"))
}
