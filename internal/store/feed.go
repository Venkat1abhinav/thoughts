package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]*PostsWithMetaData, error) {
	sort := "DESC"

	if fq.Sort == "asc" {
		sort = "ASC"
	}

	query := `
	SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username, COUNT(c.id) AS comments_count
	FROM posts p LEFT JOIN comments c ON c.post_id = p.id LEFT JOIN users u ON p.user_id = u.id LEFT JOIN followers f ON f.follower_id = p.user_id AND f.user_id = $1
	WHERE (p.user_id = $1 OR f.user_id IS NOT NULL) AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
	GROUP BY p.id, u.username
	ORDER BY p.created_at ` + sort + ` LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	rows, err := s.db.Query(
		ctx,
		query,
		userID,
		fq.Limit,
		fq.Offset,
		fq.Search,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feed, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (*PostsWithMetaData, error) {
			var post PostsWithMetaData

			err := row.Scan(
				&post.ID,
				&post.UserID,
				&post.Title,
				&post.Content,
				&post.CreatedAt,
				&post.Version,
				&post.Tags,
				&post.Username,
				&post.CommentCount,
			)
			if err != nil {
				return nil, err
			}

			return &post, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func monx(int) error {
	return nil
}
