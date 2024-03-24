package storage

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

func GetDriver() (neo4j.DriverWithContext, error) {
	dbUser := ""
	dbPassword := ""
	dbUri := "bolt://localhost:7687" // scheme://host(:port) (default port is 7687)
	return neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(dbUser, dbPassword, ""))
}
