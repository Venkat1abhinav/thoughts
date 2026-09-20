package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/owned_dragon/thoughts/gen/api"
	"github.com/owned_dragon/thoughts/internal/store"
	"github.com/rs/zerolog"
)

type application struct {
	config  config
	store   store.Store
	version string
	logger  zerolog.Logger
}

var _ api.ServerInterface = (*application)(nil)

type config struct {
	addr     string
	dbConfig dbConfig
	env      string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/docs.html")
	})

	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	})

	api.HandlerFromMux(app, r)

	return r
}

func (app *application) run(mux http.Handler) error {
	serve := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info().Str("addr", app.config.addr).Msg("server has started")
	err := serve.ListenAndServe()
	return err
}
