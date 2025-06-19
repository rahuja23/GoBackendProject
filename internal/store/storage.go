package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		GetByID(context.Context, int64) (*Post, error)
		Create(ctx context.Context, post *Post) error
	}
	Users interface {
		Create(ctx context.Context, user *User) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		GetCommentsByPostId(ctx context.Context, postId int64) ([]Comment, error)
		Delete(ctx context.Context, comment *CommentDelete) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentsStore{db: db},
	}
}
