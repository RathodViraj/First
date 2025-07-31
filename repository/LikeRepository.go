package repository

import (
	"First/model"
	"context"
	"database/sql"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type LikeRepository interface {
	AddLike(like *model.Like) error
	RemoveLike(like *model.Like) error
	changeLikeCount(postID int, delta int) error
	WhoLiked(postID int) ([]model.User, error)
}

type likeRepo struct {
	db    *sql.DB
	graph Graph
	ur    UserRepository
}

func NewLikeRepo(db *sql.DB, gr Graph, ur UserRepository) *likeRepo {
	return &likeRepo{db, gr, ur}
}

func (r *likeRepo) AddLike(like *model.Like) error {
	session := r.graph.(*graph).driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Post {id: $pid}), (u:User {id: $uid})
			MERGE (p)-[:LIKED_BY]->(u)
		`
		params := map[string]any{
			"pid": like.Pid,
			"uid": like.Uid,
		}

		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})

	if err != nil {
		return err
	}

	return r.changeLikeCount(like.Pid, 1)
}

func (r *likeRepo) RemoveLike(like *model.Like) error {
	session := r.graph.(*graph).driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (p:Post {id: $pid})-[l:LIKED_BY]->(u:User {id: $uid})
			DELETE l
		`
		params := map[string]any{
			"pid": like.Pid,
			"uid": like.Uid,
		}

		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})

	if err != nil {
		return err
	}

	// Decrease like count
	return r.changeLikeCount(like.Pid, -1)
}

func (r *likeRepo) changeLikeCount(postID int, delta int) error {
	query := `UPDATE posts SET likes = likes + ? WHERE id = ?`
	_, err := r.db.Exec(query, delta, postID)
	return err
}

func (r *likeRepo) WhoLiked(postID int) ([]model.User, error) {
	session := r.graph.(*graph).driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	ids, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (:Post {id: $postID})-[:LIKED_BY]->(u:User)
				RETURN u.id AS id`
		params := map[string]any{
			"postID": postID,
		}
		res, err := tx.Run(context.Background(), query, params)

		if err != nil {
			return nil, err
		}

		ids := []int{}
		for res.Next(context.Background()) {
			record := res.Record()
			idVal, _ := record.Get("id")
			if idInt, ok := idVal.(int64); ok {
				ids = append(ids, int(idInt))
			}
		}

		return ids, err
	})

	if err != nil {
		return nil, err
	}

	userIDs := ids.([]int)
	var users []model.User
	for id := range userIDs {
		u, err := ur.GetByID(id)
		if err == nil {
			users = append(users, *u)
		}
	}

	return users, nil
}
