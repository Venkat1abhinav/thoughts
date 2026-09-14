// Package store this is
package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound        error         = errors.New("record not found")
	QueryTimeOut       time.Duration = time.Second * 5
	ErrVersionConflict error         = errors.New("post version conflict")
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		DeleteByID(context.Context, int64) error
		UpdateByID(ctx context.Context, post *Post) (*Post, error)
		CreateMany(context.Context, []*Post) error
	}
	Users interface {
		Create(context.Context, *User) error
		CreateMany(context.Context, []*User) error
	}
	Comments interface {
		GetCommentsByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
		CreateMany(context.Context, []*Comment) error
	}
}

func NewStorage(db *pgxpool.Pool) Store {
	return Store{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
