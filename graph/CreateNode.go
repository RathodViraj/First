package graph

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Graph struct {
}

func (g *Graph) CreateUserNode(id int) error {
	session := Driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	_, err := session.ExecuteWrite(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
		query := `CREATE (u:User {id: $id})`
		params := map[string]any{"id": id}
		_, err := tx.Run(context.Background(), query, params)

		return nil, err
	})

	return err
}
