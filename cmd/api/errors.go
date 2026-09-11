package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error method: %s path: %s error: %s\n", r.Method, r.URL.Path, err)

	err = writeJSONError(w, http.StatusInternalServerError, "server encountered a problem")
	if err != nil {
		log.Println(err.Error())
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error method: %s path: %s error: %s\n", r.Method, r.URL.Path, err)

	err = writeJSONError(w, http.StatusBadRequest, err.Error())
	if err != nil {
		log.Println(err.Error())
	}
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error method: %s path: %s error: %s\n", r.Method, r.URL.Path, err)

	err = writeJSONError(w, http.StatusNotFound, err.Error())
	if err != nil {
		log.Println(err.Error())
	}
}
