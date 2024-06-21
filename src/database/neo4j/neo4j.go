package neo4j

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4j struct {
	neo4j.DriverWithContext
}

func New(config *Config) (*Neo4j, error) {
	driver, err := neo4j.NewDriverWithContext(
		config.Uri(),
		neo4j.BasicAuth(config.User, config.Pass, ""),
	)

	return &Neo4j{driver}, err
}

func (n *Neo4j) CreateEdge(parent string, child string) error {
	const query = `MERGE (parent:Article {name: $parent})
					MERGE (child:Article {name: $child})
					MERGE (parent)-[:LINK_TO]->(child)`

	ctx := context.Background()
	session := n.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"parent": parent,
		"child":  child,
	}

	_, err := session.Run(ctx, query, params)
	return err
}
