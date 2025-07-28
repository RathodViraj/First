package repository

import (
	"First/model"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ConnectionRepository interface {
	CreateConnection(conn *model.Connection) error
	DeleteConnection(followerID, followingID int) error
	GetFollowers(userIDs []int) []model.User
	GetFollowings(userIDs []int) []model.User
	GetMutual(mutualIDs []int) ([]model.User, error)
}

type connectionRepo struct {
	db    *sql.DB
	graph Graph
}

func NewConnectionRepo(db *sql.DB, graph Graph) *connectionRepo {
	return &connectionRepo{db: db, graph: graph}
}

func (r *connectionRepo) CreateConnection(conn *model.Connection) error {
	if conn.FollowerID == conn.FollowingID {
		return fmt.Errorf("cannot follow yourself")
	}

	session := r.graph.(*graph).driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (follower:User {id: $followerID}), (following:User {id: $followingID})
		CREATE (follower)-[:FOLLOWS]->(following)`
		params := map[string]any{
			"followerID":  conn.FollowerID,
			"followingID": conn.FollowingID,
		}
		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})
	return err
}

func (r *connectionRepo) DeleteConnection(followerID, followingID int) error {
	session := r.graph.(*graph).driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (follower:User {id: $followerID})-[r:FOLLOWS]->(following:User {id: $followingID}) DELETE r`
		params := map[string]any{
			"followerID":  followerID,
			"followingID": followingID,
		}
		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})

	return err
}

func (r *connectionRepo) GetFollowers(userIDs []int) []model.User {
	var followers []model.User

	for _, id := range userIDs {
		query := `SELECT id, name, email, created_at FROM users WHERE id = ?`
		row := r.db.QueryRow(query, id)

		var u model.User
		if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			log.Println(err)
			continue
		}
		followers = append(followers, u)
	}

	return followers
}

func (r *connectionRepo) GetFollowings(userIDs []int) []model.User {
	var followings []model.User

	for _, id := range userIDs {
		query := `SELECT id, name, email, created_at FROM users WHERE id = ?`
		row := r.db.QueryRow(query, id)

		var u model.User
		if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			log.Println(err)
			continue
		}
		followings = append(followings, u)
	}

	return followings
}

func (r *connectionRepo) GetMutual(mutualIDs []int) ([]model.User, error) {
	var mutuals []model.User

	for _, id := range mutualIDs {
		query := `SELECT id, name, email, created_at FROM users WHERE id = ?`
		row := r.db.QueryRow(query, id)

		var u model.User
		if err := row.Scan(&u.Id, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			log.Println(err)
			continue
		}
		mutuals = append(mutuals, u)
	}

	return mutuals, nil
}
