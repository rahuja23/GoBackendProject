package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	Tags      []string `json:"tags"`

	UpdatedAt string `json:"updated_at"`

	Version  int64            `json:"version"`
	Comments []Comment        `json:"comments"`
	User     PostUserMetadata `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
}
type DeletePost struct {
	ID int64 `json:"id"`
}
type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
	SET title= $1, content = $2,  version = version +1
	WHERE id = $3 AND version = $4
	RETURNING version`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil

}
func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags, version)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
		0,
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
func (s *PostsStore) Get(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, user_id, title, content,  created_at, updated_at, version,  tags  FROM posts
	where id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
		pq.Array(&post.Tags),
	)

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

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts where id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userId int64) ([]PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at,   p.version, p.tags, 
			u.username, u.id, u.email,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN users u ON p.user_id= u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id=$1
		WHERE f.user_id =$1 OR p.user_id=$1
		GROUP BY p.id, u.username, u.id
		ORDER BY p.created_at DESC;
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err = rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.User.ID,
			&p.User.Email,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}
	return feed, nil
}
