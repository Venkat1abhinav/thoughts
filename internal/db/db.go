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
	log.Println("database connected")
	maxIdleDurationm, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(maxOpenConns)
	config.MinConns = int32(maxIdleConns)
	config.MaxConnIdleTime = maxIdleDurationm

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
