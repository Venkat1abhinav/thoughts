package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/owned_dragon/thoughts/internal/store"
)

type application struct {
	config  config
	store   store.Store
	version string
}

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
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Get("/", app.healthCheckHandler())
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler())
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler())
			r.Route("/{id}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler())
				r.Delete("/", app.deletePostHandler())
				r.Patch("/", app.updatePostHandler())
			})
		})
	})
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

	log.Printf("starting the server at %s", app.config.addr)

	err := serve.ListenAndServe()
	return err
}
