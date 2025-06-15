package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
	CreatedAt string `json:"created_at"`
	USER      User   `json:"user"`
}
type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
	INSERT INTO comments (post_id,  user_id,  content)
	VALUES ($1, $2, $3) RETURNING id, created_at
	`
	err := s.db.QueryRowContext(
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

func (s *CommentsStore) GetCommentsByPostId(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, c.created_at 
	FROM comments c 
	JOIN users on c.user_id = users.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC;
	`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.USER = User{}
		fmt.Println(rows)
		err = rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)

	}
	return comments, nil
}
