package storage

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateOpening(openingName string) error {
	driver, err := GetDriver()
	if err != nil {
		return err
	}

	ctx := context.Background()
	defer driver.Close(ctx)

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return err
	}

	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer session.Close(ctx)

	indexes := []string{
		"CREATE INDEX ON :Opening(name);",
	}

	constraints := []string{
		"CREATE CONSTRAINT ON (o:Opening) ASSERT o.name IS UNIQUE;",
	}

	// Run index queries via implicit auto-commit transaction
	for _, index := range indexes {
		_, err = session.Run(ctx, index, nil)
		if err != nil {
			return err
		}
	}

	// Run constraint queries via implicit auto-commit transaction
	for _, constraint := range constraints {
		_, err = session.Run(ctx, constraint, nil)
		if err != nil {
			return err
		}

	}

	query := fmt.Sprintf("CREATE (o:Opening {name: \"%s\"})", openingName)

	_, err = neo4j.ExecuteQuery(ctx, driver, query, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return err
	}

	return nil
}

func GetOpenings() ([]string, error) {
	driver, err := GetDriver()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	defer driver.Close(ctx)

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}

	query := "MATCH (n:Opening{}) RETURN n;"
	result, err := neo4j.ExecuteQuery(ctx, driver, query, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	var openings []string
	for _, node := range result.Records {
		opening := node.AsMap()["n"].(neo4j.Node).Props["name"].(string)
		openings = append(openings, opening)
	}

	return openings, nil
}
