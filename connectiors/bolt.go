package connectiors

import (
	"fmt"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/pkg/errors"
	"gitlab.com/build-graph/models"
)

type BoltConnection struct {
	Connection bolt.Conn
}

func GetBoltConnection(cfg *models.Config) (*BoltConnection, error) {

	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(fmt.Sprintf("bolt://%s:%s@%s:7687", cfg.BoltUser, cfg.BoltPassword, cfg.BoltHost))
	if err != nil {
		return nil, errors.New("Neo4j Connection not established .")
	}
	//defer conn.Close()

	return &BoltConnection{conn}, nil
}
