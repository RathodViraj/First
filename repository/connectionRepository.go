package repository

import (
	"First/model"
	"database/sql"
	"fmt"
)

type ConnectionRepository interface {
	CreateConnection(conn *model.Connection) error
	DeleteConnection(followerID, followingID int) error
	GetFollowers(userID int) ([]model.User, error)
	GetFollowings(userID int) ([]model.User, error)
	GetMutual(id int) ([]model.User, error)
}

type connectionRepo struct {
	db *sql.DB
}

func NewConnectioRepo(db *sql.DB) *connectionRepo {
	return &connectionRepo{db}
}

func (r *connectionRepo) CreateConnection(conn *model.Connection) error {
	if conn.FollowerID == conn.FollowingID {
		return fmt.Errorf("cannot follow yourself")
	}

	query := `
        INSERT INTO connections (follower_id, following_id) 
        VALUES (?, ?)`
	_, err := r.db.Exec(query, conn.FollowerID, conn.FollowingID)
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}
	return nil
}

func (r *connectionRepo) DeleteConnection(followerID, followingID int) error {
	res, err := r.db.Exec(
		"DELETE FROM connections WHERE follower_id = ? AND following_id = ?",
		followerID, followingID,
	)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("connection not found")
	}
	return nil
}

func (r *connectionRepo) GetFollowers(userID int) ([]model.User, error) {
	query := `SELECT u.id, u.name, u.email, u.created_at 
        FROM users u
        JOIN connections c ON u.id = c.follower_id
        WHERE c.following_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var followers []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		followers = append(followers, u)
	}

	return followers, nil
}

func (r *connectionRepo) GetFollowings(userID int) ([]model.User, error) {
	query := `
        SELECT u.id, u.name, u.email, u.created_at 
        FROM users u
        JOIN connections c ON u.id = c.following_id
        WHERE c.follower_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var followings []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		followings = append(followings, u)
	}
	return followings, nil
}

func (r *connectionRepo) GetMutual(userID int) ([]model.User, error) {
	query := `
		SELECT u.id, u.name, COUNT(*) AS mutual_count
		FROM users u
		JOIN connections c1 ON u.id = c1.following_id
		JOIN connections c2 ON c1.following_id = c2.follower_id
		WHERE c1.follower_id = ? AND c2.following_id = ? AND u.id != ?
		GROUP BY u.id, u.name
		ORDER BY mutual_count DESC
	`

	rows, err := r.db.Query(query, userID, userID, userID)
	if err != nil {
		return []model.User{}, nil
	}

	var mutual []model.User
	for rows.Next() {
		var u model.User
		if err = rows.Scan(&u.Id, &u.Name); err != nil {
			return []model.User{}, nil
		}

		mutual = append(mutual, u)
	}

	return mutual, nil
}
