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
		SELECT id, content, title, user_id, tags, created_at, updated_at
		FROM posts
		WHERE id = $1
	`
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

	_, err := s.db.Exec(ctx, query, postID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) UpdateByID(
	ctx context.Context,
	title string,
	content string,
	postID int64,
) (*Post, error) {
	var post Post
	query := `
		UPDATE posts
		SET
			title = $1,
			content = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING id, title, content, user_id, tags, created_at, updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		title,
		content,
		postID,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
