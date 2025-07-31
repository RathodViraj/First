package repository

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type graph struct {
	driver neo4j.DriverWithContext
}

type Graph interface {
	CreateUserNode(id int) error
	DeleteUserNode(id int) error
	GetFollowersIDs(userID int) ([]int, error)
	GetFollowingsIDs(userID int) ([]int, error)
	GetMutualIDs(userID int) ([]int, error)
	CreatePostNode(id, userID, parentID int) error
	DeletePostNode(id int) error
	GetUserIDByPostID(id int) (int, error)
}

func NewGraph(d neo4j.DriverWithContext) Graph {
	return &graph{d}
}

func (g *graph) CreateUserNode(id int) error {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MERGE (u:User {id: $id})`
		params := map[string]any{"id": id}
		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})

	return err
}

func (g *graph) DeleteUserNode(id int) error {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {id: $id}) DELETE u`
		params := map[string]any{"id": id}
		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})

	return err
}

func (g *graph) GetFollowersIDs(userID int) ([]int, error) {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `MATCH (:User {id: $userID})<-[:FOLLOWS]-(u:User)
				RETURN u.id AS id`
		params := map[string]any{"userID": userID}

		res, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}

		var followerIDs []int
		for res.Next(context.Background()) {
			record := res.Record()
			idVal, _ := record.Get("id")
			if idInt, ok := idVal.(int64); ok {
				followerIDs = append(followerIDs, int(idInt))
			}
		}

		return followerIDs, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]int), nil
}

func (g *graph) GetFollowingsIDs(userID int) ([]int, error) {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `MATCH (:User {id: $userID})-[:FOLLOWS]->(u:User)
				RETURN u.id AS id`
		params := map[string]any{"userID": userID}

		res, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}

		var followingIDs []int
		for res.Next(context.Background()) {
			record := res.Record()
			idVal, _ := record.Get("id")
			if idInt, ok := idVal.(int64); ok {
				followingIDs = append(followingIDs, int(idInt))
			}
		}

		return followingIDs, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]int), nil
}

func (g *graph) GetMutualIDs(userID int) ([]int, error) {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `MATCH (u:User {id: $userID})-[:FOLLOWS]->(mutual:User)-[:FOLLOWS]->(u)
				RETURN mutual.id AS id`
		params := map[string]any{"userID": userID}

		res, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}

		var mutualIDs []int
		for res.Next(context.Background()) {
			record := res.Record()
			idVal, _ := record.Get("id")
			if idInt, ok := idVal.(int64); ok {
				mutualIDs = append(mutualIDs, int(idInt))
			}
		}

		return mutualIDs, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]int), nil
}

func (g *graph) CreatePostNode(id, userID, parentID int) error {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (u:User {id: $userID}) CREATE (p:Post {id: $id, pid: $parentID})-[:POSTED_BY]->(u)`
		params := map[string]any{
			"id":       id,
			"userID":   userID,
			"parentID": parentID,
		}
		_, err := tx.Run(context.Background(), query, params)
		return nil, err
	})
	return err
}

func (g *graph) DeletePostNode(id int) error {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (p:Post {id: $id}) DETACH DELETE p`
		params := map[string]any{"id": id}
		_, err := tx.Run(context.Background(), query, params)

		return nil, err
	})

	return err
}

func (g *graph) GetUserIDByPostID(id int) (int, error) {
	session := g.driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	// Get the user ID who posted a given post
	res, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MATCH (p:Post {id: $id})-[:POSTED_BY]->(u:User) RETURN u.id AS userId`
		params := map[string]any{"id": id}
		result, err := tx.Run(context.Background(), query, params)
		if err != nil {
			return nil, err
		}
		if result.Next(context.Background()) {
			record := result.Record()
			idVal, _ := record.Get("userId")
			if idInt, ok := idVal.(int64); ok {
				return int(idInt), nil
			}
		}
		return nil, result.Err()
	})
	if err != nil {
		return -1, err
	}
	if userID, ok := res.(int); ok {
		return userID, nil
	}
	return -1, nil
}
