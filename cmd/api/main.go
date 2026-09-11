package main

import (
	"log"

	"github.com/owned_dragon/thoughts/internal/db"
	"github.com/owned_dragon/thoughts/internal/env"
	"github.com/owned_dragon/thoughts/internal/store"
)

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

	db, err := db.New(
		cfg.dbConfig.addr,
		cfg.dbConfig.maxOpenConns,
		cfg.dbConfig.maxIdleConns,
		cfg.dbConfig.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		store:   store,
		version: "0.0.1",
	}
	mux := app.mount()

	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
