package repository

import (
	"First/model"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type UserRepository interface {
	GetByEmail(email string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	CreateUserSQL(user *model.User) error
	Delete(id, userID int, isAdmin bool) error
	DeleteUserSQL(id int) error
	GetUserFeed(userID, offset int, graph Graph) ([]model.Post, error)
	Update(user *model.User) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	query := `SELECT id, name, email, password, created_at 
              FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func (r *userRepo) CreateUserSQL(user *model.User) error {
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

func (r *userRepo) DeleteUserSQL(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepo) Delete(id, userID int, isAdmin bool) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if id != userID && !isAdmin {
		return fmt.Errorf("Sorry your not an Admin.\nOnly Admin can delete other's account.")
	}

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

func (r *userRepo) GetUserFeed(userID, offset int, graph Graph) ([]model.Post, error) {
	followingIDs, err := graph.GetFollowingsIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("neo4j error: %w", err)
	}

	if len(followingIDs) == 0 {
		return []model.Post{}, nil
	}

	query := `
	SELECT p.id, p.uid, p.content, p.likes, p.parent_id, p.created_at
	FROM posts p
	WHERE p.uid IN (?
	` + strings.Repeat(",?", len(followingIDs)-1) + `)
	ORDER BY p.created_at DESC
	LIMIT 10 OFFSET ?
	`

	args := make([]any, len(followingIDs)+1)
	for i, id := range followingIDs {
		args[i] = id
	}
	args[len(args)-1] = offset

	rows, err := r.db.Query(query, args...)
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

func (r *userRepo) Update(user *model.User) error {
	query := `UPDATE users SET name = ?, email = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
