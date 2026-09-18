package main

import (
	"log"
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check if the API is healthy
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/v1/health [get]
func (app *application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.version,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		err := writeJSONError(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func (app *application) HealthCheckRoot(
	w http.ResponseWriter,
	r *http.Request,
) {
	app.HealthCheck(w, r)
}
