package repository

import (
	"First/model"
	"database/sql"
)

type LikeRepository interface {
	AddLike(like *model.Like) error
	RemoveLike(like *model.Like) error
	changeLikeCount(postID int, delta int) error
	WhoLiked(postID int) ([]model.User, error)
}

type likeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) *likeRepo {
	return &likeRepo{db}
}

func (r *likeRepo) AddLike(like *model.Like) error {
	query := `INSERT INTO likes (pid, uid) VALUES (?, ?)`
	if _, err := r.db.Exec(query, like.Pid, like.Uid); err != nil {
		return err
	}
	return r.changeLikeCount(like.Pid, 1)
}

func (r *likeRepo) RemoveLike(like *model.Like) error {
	query := `DELETE FROM likes WHERE pid = ? AND uid = ?`
	if _, err := r.db.Exec(query, like.Pid, like.Uid); err != nil {
		return err
	}
	return r.changeLikeCount(like.Pid, -1)
}

func (r *likeRepo) changeLikeCount(postID int, delta int) error {
	query := `UPDATE posts SET likes = likes + ? WHERE id = ?`
	_, err := r.db.Exec(query, delta, postID)
	return err
}

func (r *likeRepo) WhoLiked(postID int) ([]model.User, error) {
	query := `
		SELECT u.id, u.name
		FROM users u
		JOIN likes l ON u.id = l.uid
		WHERE l.pid = ?`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
