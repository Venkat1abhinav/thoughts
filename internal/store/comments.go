package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"cotent"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `json:"user"`
}

type CommentsStore struct {
	db *pgxpool.Pool
}

func (s *CommentsStore) GetCommentsByPostID(
	ctx context.Context,
	postID int64,
) ([]Comment, error) {
	query := `
	SELECT
		c.id AS comment_id,
		c.post_id,
		c.user_id,
		c.content,
		c.created_at,
		u.username
	FROM comments AS c
	JOIN users AS u ON u.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	rows, err := s.db.Query(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
	INSERT into comments (post_id, user_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)

	defer cancel()

	err := s.db.QueryRow(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentsStore) CreateMany(
	ctx context.Context,
	comments []*Comment,
) error {
	if len(comments) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	_, err := s.db.CopyFrom(
		ctx,
		pgx.Identifier{"comments"},
		[]string{
			"post_id",
			"user_id",
			"content",
		},
		pgx.CopyFromSlice(len(comments), func(i int) ([]any, error) {
			c := comments[i]

			return []any{
				c.PostID,
				c.UserID,
				c.Content,
			}, nil
		}),
	)

	return err
}
