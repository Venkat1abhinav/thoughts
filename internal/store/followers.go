package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FollwersStore struct {
	db *pgxpool.Pool
}

// type followers struct {
// 	UserID     int64     `json:"user_id"`
// 	FollowerID int64     `json:"follower_id"`
// 	CreatedAt  time.Time `json:"created_at"`
// }

func (s *FollwersStore) Follow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	INSERT INTO FOLLOWERS (user_id, follower_id)
	Values ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrConfict
		}
		return err
	}
	return nil
}

func (s *FollwersStore) UnFollow(ctx context.Context, followerID int64, userID int64) error {
	query := `
	DELETE from followers
	where user_id = $1 and follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := s.db.Exec(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}
