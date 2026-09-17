package main

import (
	"log"

	"github.com/owned_dragon/thoughts/internal/db"
	"github.com/owned_dragon/thoughts/internal/env"
	"github.com/owned_dragon/thoughts/internal/store"
)

func main() {
	dsn := "postgres://admin:admin123@localhost:5432/thoughts?sslmode=disable"

	addr := env.GetString("DB_ADDR", dsn)

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	store := store.NewStorage(conn)

	defer conn.Close()

	err = db.Seed(store)
	if err != nil {
		log.Fatal(err)
	}
}
