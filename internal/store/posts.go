package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}
type DeletePost struct {
	ID int64 `json:"id"`
}
type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) UpdateByID(ctx context.Context, post *Post) (*Post, error) {
	query1 := `
	UPDATE posts 
	SET content = $1, 
		updated_at = NOW()
	WHERE id = $2;
	`
	query2 := `SELECT id, user_id, title, content,  created_at, updated_at, tags 
	FROM posts
	where id = $1;`
	var post_out Post
	_, err := s.db.QueryContext(ctx, query1, post.Content, post.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	err = s.db.QueryRowContext(ctx, query2, post.ID).Scan(
		&post_out.ID,
		&post_out.UserID,
		&post_out.Title,
		&post_out.Content,
		&post_out.CreatedAt,
		&post_out.UpdatedAt,
		pq.Array(&post_out.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post_out, nil

}
func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, user_id, title, content,  created_at, updated_at, tags FROM posts
	where id = $1
	`

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Content,
		&post.Title,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil

}

func (s *PostsStore) Delete(ctx context.Context, postdel *DeletePost) error {
	query := `
		DELETE FROM posts where id = $1
`
	_, err := s.db.QueryContext(ctx, query, postdel.ID)
	if err != nil {
		return err
	}
	return nil
}
