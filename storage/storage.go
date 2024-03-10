package storage

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/notnil/chess"
)

func StoreGame(game *chess.Game) {
	dbUser := ""
	dbPassword := ""
	dbUri := "bolt://localhost:7687" // scheme://host(:port) (default port is 7687)
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
	ctx := context.Background()
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Viola! Connected to Memgraph!")
	}

	//Create a simple session
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer session.Close(ctx)

	/*
		// Run index queries via implicit auto-commit transaction
		for _, index := range indexes {
			_, err = session.Run(ctx, index, nil)
			if err != nil {
				panic(err)
			}
		}
		fmt.Println("****** Indexes created *******")*/

	indexes := []string{
		"CREATE INDEX ON :Position(fen);",
	}

	// Initial position first
	query := "MERGE (n0:Position {fen: \"start_position\"})"

	// Create all nodes from the game history
	for i, history := range game.MoveHistory() {
		if history.PostPosition == nil {
			// Already recorded as we record the next position always
			break
		}

		// Create a node for the position
		query += fmt.Sprintf("\nMERGE (n%d)-[:NEXT]->(n%d:Position {fen: \"%s\"})", i, i+1, history.PostPosition.String())
	}

	fmt.Println(query)

	// Run index queries via implicit auto-commit transaction
	for _, index := range indexes {
		_, err = session.Run(ctx, index, nil)
		if err != nil {
			panic(err)
		}
	}

	_, err = neo4j.ExecuteQuery(ctx, driver, query, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		panic(err)
	}
}
