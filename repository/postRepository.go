package repository

import (
	"First/model"
	"database/sql"
	"errors"
	"fmt"
)

type PostRepository interface {
	Create(post *model.Post) error
	Delete(id, uid int, isAdmin bool) error
	GetByID(id int) (*model.Post, error)
	GetRecentPosts(offset int) ([]model.Post, error)
	GetAllUserPosts(id int) (*[]model.Post, error)
}

var ur userRepo
var cr commentRepo

type postRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) PostRepository {
	return &postRepo{db}
}

func (r *postRepo) Create(post *model.Post) error {
	query := `
        INSERT INTO posts (uid, content, likes) 
        VALUES (?, ?, 0)`

	res, err := r.db.Exec(query, post.Uid, post.Content, post.Likes)
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

func (r *postRepo) Delete(postID, userID int, isAdmin bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete comments
	if _, err := tx.Exec("DELETE FROM posts WHERE parent_id = ?", postID); err != nil {
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	var res sql.Result
	if isAdmin {
		res, err = tx.Exec("DELETE FROM posts WHERE id = ?", postID)
	} else {
		res, err = tx.Exec("DELETE FROM posts WHERE id = ? AND uid = ?", postID, userID)
	}
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("post not found or not authorized")
	}

	return tx.Commit()
}

func (r *postRepo) GetByID(id int) (*model.Post, error) {
	query := `
        SELECT id, uid, content, likes, parent_id, created_at
        FROM posts WHERE id = ?`

	var post model.Post
	err := r.db.QueryRow(query, id).Scan(
		&post.Id, &post.Uid, &post.Content,
		&post.Likes, &post.ParentId, &post.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	comments, err := cr.GetComments(id)
	if err != nil {
		return nil, err
	}

	post.Comments = comments

	return &post, nil
}

func (r *postRepo) GetRecentPosts(offset int) ([]model.Post, error) {
	query := `
        SELECT id, uid, content, likes, parent_id, created_at
        FROM posts 
        WHERE parent_id IS NULL
        ORDER BY created_at DESC
        LIMIT 10 OFFSET ?`

	rows, err := r.db.Query(query, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent posts: %w", err)
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.Id, &post.Uid, &post.Content,
			&post.Likes, &post.ParentId, &post.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *postRepo) GetAllUserPosts(id int) (*[]model.Post, error) {
	var posts []model.Post

	query := `
        SELECT id, uid, content, likes, created_at
        FROM posts WHERE uid = ?
	    AND parent_id IS NULL`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("can't fetch posts")
	}

	for rows.Next() {
		var post model.Post
		if err := rows.Scan(
			&post.Id, &post.Uid, &post.Content,
			&post.Likes, &post.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		posts = append(posts, post)
	}

	return &posts, nil
}
