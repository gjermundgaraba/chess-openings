package storage

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/notnil/chess"
)

type MovesAtPositions map[string][]string

func StoreVariation(game *chess.Game, opening string) {
	driver, err := GetDriver()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	//Create a simple session
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""})
	defer session.Close(ctx)

	indexes := []string{
		"CREATE INDEX ON :Position(fen);",
		"CREATE INDEX ON :Position(moves_made);",
	}

	// Initial position first
	query := fmt.Sprintf(`MATCH (o:Opening {name: "%s"})
MERGE (n0:Position {fen: "%s", moves_made: 0})
MERGE (n0)-[:IN]->(o)`, opening, startFEN)

	// Create all nodes from the game history
	for i, history := range game.MoveHistory() {
		if history.PostPosition == nil {
			// Already recorded as we record the next position always
			break
		}

		// Create a node for the position
		notation := chess.AlgebraicNotation{}.Encode(history.PrePosition, history.Move)
		query += fmt.Sprintf(`
MERGE (n%d)-[:NEXT {move: "%s"}]->(n%d:Position {fen: "%s", moves_made: %d})
MERGE (n%d)-[:IN]->(o)`, i, notation, i+1, history.PostPosition.String(), i+1, i+1)
	}

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

func GetVariations(opening string) (MovesAtPositions, error) {
	driver, err := GetDriver()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	defer driver.Close(ctx)

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`MATCH (p:Position)-[:IN]->(:Opening{name: "%s"})
MATCH (p)-[next:NEXT]->(:Position)
RETURN p, next
ORDER BY p.moves_made`, opening)

	result, err := neo4j.ExecuteQuery(ctx, driver, query, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	movesInOpening := make(MovesAtPositions)
	for _, record := range result.Records {
		node := record.AsMap()["p"].(neo4j.Node)
		fen := node.Props["fen"].(string)
		moveStr := record.AsMap()["next"].(neo4j.Relationship).Props["move"].(string)
		movesInOpening[fen] = append(movesInOpening[fen], moveStr)
	}

	return movesInOpening, nil
}
