// Package db covers all database layer
package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(addr string,
	maxOpenConns, maxIdleConns int,
	maxIdleTime string,
) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()
	db, err := pgxpool.New(
		ctx,
		addr,
	)
	if err != nil {
		return nil, err
	}

	log.Println("database connected")
	maxIdleDurationm, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	config := db.Config()
	config.MaxConns = int32(maxOpenConns)
	config.MinConns = int32(maxIdleConns)
	config.MaxConnIdleTime = maxIdleDurationm

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
