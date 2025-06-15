package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s and error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, "An internal server error has occurred")

}

func (app *application) badrequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s and error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())

}

func (app *application) notfoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found  error: %s path: %s and error: %s", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, "Resource not found ")

}
