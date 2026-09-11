package store

import (
	"context"
	"time"

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

func (s *CommentsStore) GetPostByID(
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
