package repository

import (
	"First/model"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type CommentsRepository interface {
	GetComments(id int) ([]model.Post, error)
	AddComment(post *model.Post) error
}

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *commentRepo {
	return &commentRepo{db}
}

func (r *commentRepo) GetComments(id int) ([]model.Post, error) {
	commentsQuery := `
        SELECT id, uid, content, likes, parent_id, created_at
        FROM posts WHERE parent_id = ?
        ORDER BY created_at ASC`

	rows, err := r.db.Query(commentsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Post
	for rows.Next() {
		var comment model.Post
		if err := rows.Scan(
			&comment.Id, &comment.Uid, &comment.Content,
			&comment.Likes, &comment.ParentId, &comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *commentRepo) AddComment(post *model.Post) error {
	if strings.TrimSpace(post.Content) == "" {
		return errors.New("comment content cannot be empty")
	}
	if post.ParentId == nil {
		return errors.New("parent_id is required for a comment")
	}

	query := `
	INSERT INTO posts (uid, content, likes, parent_id) 
	VALUES (?, ?, 0, ?)`

	res, err := r.db.Exec(query, post.Uid, post.Content, post.Likes, post.ParentId)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	post.Id = int(id)
	return nil
}
