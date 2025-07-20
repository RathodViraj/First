package repository

import (
	"First/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Create(user *model.User) error
	Delete(id, userID int, isAdmin bool) error
	GetUserFeed(id, offset int) ([]model.Post, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	// Normalize email to lowercase and trim spaces
	email = strings.TrimSpace(strings.ToLower(email))

	query := `SELECT id, name, email, password, created_at 
              FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Log the email that wasn't found for debugging
			log.Printf("User not found for email: %s", email)
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

func (r *userRepo) GetByID(id int) (*model.User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Create(user *model.User) error {
	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	res, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	user.Id = int(id)
	return nil
}

func (r *userRepo) Delete(id, userID int, isAdmin bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Safe to call if tx.Commit() succeeds

	if id != userID && !isAdmin {
		return fmt.Errorf("Sorry your not an Admin.\nOnly Admin can delete other's account.")
	}

	// Delete connections first to avoid FK violations
	if _, err = tx.Exec("DELETE FROM connections WHERE follower_id = ? OR following_id = ?", id, id); err != nil {
		return fmt.Errorf("failed to delete user connections: %w", err)
	}

	res, err := tx.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user not found")
	}

	return tx.Commit()
}

func (r *userRepo) GetUserFeed(id, offset int) ([]model.Post, error) {
	query := `
	SELECT p.id, p.uid, p.content, p.likes, p.parent_id, p.created_at
	FROM posts p
	JOIN connections c ON p.uid = c.following_id
	WHERE c.follower_id = ? 
	ORDER BY p.created_at DESC
	LIMIT 10 OFFSET ?
	`
	rows, err := r.db.Query(query, id, offset)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var feed []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.Id, &p.Uid, &p.Content, &p.Likes, &p.ParentId, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		feed = append(feed, p)
	}
	return feed, nil
}
