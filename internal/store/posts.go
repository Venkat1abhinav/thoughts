package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostGet struct {
	Content string    `json:"content"`
	Title   string    `json:"title"`
	UserID  int64     `json:"user_id"`
	Tags    []string  `json:"tags"`
	Comment []Comment `json:"comments"`
}

type PostsStore struct {
	db *pgxpool.Pool
}

func (s *PostsStore) Create(ctx context.Context, p *Post) error {
	query := `
		INSERT into posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		p.Content,
		p.Title,
		p.UserID,
		p.Tags,
	).Scan(
		&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, version,created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	var post Post

	err := s.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		&post.Tags,
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err

		}
	}

	return &post, nil
}

func (s *PostsStore) DeleteByID(
	ctx context.Context, postID int64,
) error {
	query := `DELETE from posts where id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	res, err := s.db.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	rows := res.RowsAffected()

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostsStore) UpdateByID(
	ctx context.Context,
	p *Post,
) (*Post, error) {
	var post Post

	query := `
        UPDATE posts
        SET
            title = CASE
                WHEN $1::text IS NULL THEN title
                ELSE $1
            END,
            content = CASE
                WHEN $2::text IS NULL THEN content
                ELSE $2
            END,
			version = version + 1,
            updated_at = NOW()
        WHERE id = $3 AND version = $4
        RETURNING id, title, content, user_id, tags, version, created_at, updated_at;
    `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		p.Title,
		p.Content,
		p.ID,
		p.Version,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Tags,
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
