package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersStore struct {
	db *pgxpool.Pool
}

func (s *UsersStore) Create(ctx context.Context, u *User) error {
	query := `
		INSERT into users (username, first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		u.Username,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
	).Scan(
		&u.ID,
		&u.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
