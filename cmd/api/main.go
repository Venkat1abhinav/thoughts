package main

import (
	"os"

	"github.com/owned_dragon/thoughts/internal/db"
	"github.com/owned_dragon/thoughts/internal/env"
	"github.com/owned_dragon/thoughts/internal/store"
	"github.com/rs/zerolog"
)

//	@title			thoughts API
//	@version		1.0
//	@description	Rest API for Thoughs a Distraction Free Social Application
//	@host			localhost:3000
//	@BasePath		/

func main() {
	dsn := "postgres://admin:admin123@localhost:5432/thoughts?sslmode=disable"
	dbConfig := dbConfig{
		addr: env.GetString(
			"DB_ADDR",
			dsn,
		),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	cfg := config{
		addr:     env.GetString("ADDR", ":8000"),
		dbConfig: dbConfig,
		env:      env.GetString("ENV", "development"),
	}

	var logger zerolog.Logger
	if cfg.env == "development" {
		logger = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Logger()
	} else {
		logger = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Logger()
	}

	db, err := db.New(
		cfg.dbConfig.addr,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Failed to connect to the database")
	}

	logger.Info().Msg("connected to the database")

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		store:   store,
		version: "0.0.1",
		logger:  logger,
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
