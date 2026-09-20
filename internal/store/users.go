package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hashed
	return nil
}

type UsersStore struct {
	db *pgxpool.Pool
}

func (s *UsersStore) Create(ctx context.Context, u *User) error {
	query := `
	Insert into users (username, first_name, last_name, email, password)
	Values ($1, $2, $3, $4, $5) RETURNING id, created_at
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

func (s *UsersStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
	SELECT id, username, first_name, last_name, email, password, created_at
	from users
	WHERE id=$1
	`

	var user User
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UsersStore) CreateMany(
	ctx context.Context,
	users []*User,
) error {
	if len(users) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	args := make([]any, 0, len(users)*5)
	values := make([]string, 0, len(users))

	for i, u := range users {
		offset := i * 5

		values = append(values, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d)",
			offset+1,
			offset+2,
			offset+3,
			offset+4,
			offset+5,
		))

		args = append(
			args,
			u.Username,
			u.FirstName,
			u.LastName,
			u.Email,
			u.Password,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO users (
			username,
			first_name,
			last_name,
			email,
			password
		)
		VALUES %s
		RETURNING id, created_at
	`, strings.Join(values, ", "))

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for i := range users {
		if !rows.Next() {
			return rows.Err()
		}

		if err := rows.Scan(
			&users[i].ID,
			&users[i].CreatedAt,
		); err != nil {
			return err
		}
	}

	return rows.Err()
}
