package repository

import (
	"First/model"
	"database/sql"
	"log"
)

type searchRepo struct {
	db *sql.DB
}

type SearchRepository interface {
	SearchUser(key string) ([]model.User, error)
	SearchPost(key string) ([]model.Post, error)
}

func NewSearchRepo(db *sql.DB) SearchRepository {
	return &searchRepo{db}
}

func (r *searchRepo) SearchUser(key string) ([]model.User, error) {
	query := `SELECT ID, NAME, ROLE FROM USERS WHERE NAME LIKE ?`

	rows, err := r.db.Query(query, "%"+key+"%")
	if err != nil {
		log.Println(err.Error())
		return []model.User{}, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err = rows.Scan(&u.Id, &u.Name, &u.Role); err == nil {
			users = append(users, u)
		}
	}

	return users, nil
}

func (r *searchRepo) SearchPost(key string) ([]model.Post, error) {
	query := `SELECT ID, UID, CONTENT, LIKES FROM POSTS WHERE CONTENT LIKE ?`

	rows, err := r.db.Query(query, "%"+key+"%")
	if err != nil {
		return []model.Post{}, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err = rows.Scan(&p.Id, &p.Uid, &p.Content, &p.Likes); err == nil {
			posts = append(posts, p)
		}
	}

	return posts, nil
}
