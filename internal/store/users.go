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
	Insert into users (username, first_name, last_name, email, password, created_at)
	Values ($1, $2, $3, $4, $5, $6) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

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
