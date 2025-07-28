package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource conflict")
	QueryTimeoutDuration = time.Minute * 2
)

type Storage struct {
	Posts interface {
		Get(context.Context, int64) (*Post, error)
		Create(ctx context.Context, post *Post) error
		Delete(ctx context.Context, id int64) error
		Update(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, userId int64) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		Get(ctx context.Context, id int64) (*User, error)
		Update(ctx context.Context, user *User) (*User, error)
		Delete(ctx context.Context, id int64) error
	}
	Comments interface {
		Create(ctx context.Context, comment *Comment) error
		Get(ctx context.Context, postId int64) ([]Comment, error)
		Delete(ctx context.Context, comment *CommentDelete) error
		Update(ctx context.Context, comment *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context, followerID, userID int64) error
		Unfollow(ctx context.Context, followerID, userID int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		&PostsStore{db: db},
		&UserStore{db: db},
		&CommentsStore{db: db},
		&FollowerStore{db: db},
	}
}
