package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Err(err).
		Msg("internal server error")
	err = writeJSONError(w, http.StatusInternalServerError, "server encountered a problem")
	app.logger.Err(err)

	if err != nil {
		app.logger.Info().Err(err).Msg(err.Error())
	}
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Err(err).
		Msg("bad request error")

	err = writeJSONError(w, http.StatusBadRequest, err.Error())
	app.logger.Err(err)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Err(err).
		Msg("not found error")

	err = writeJSONError(w, http.StatusNotFound, err.Error())
	app.logger.Err(err)
}

func (app *application) confictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Info().Str("method", r.Method).Str("path", r.URL.Path).Err(err).Msg("conflict error")

	err = writeJSONError(w, http.StatusConflict, err.Error())
	app.logger.Err(err)
}
