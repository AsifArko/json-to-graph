package main

import (
	json "encoding/json"
	"fmt"
	"github.com/couchbase/gocb"
	J "github.com/restra-social/jypher"
	jmodel "github.com/restra-social/jypher/models"
	"gitlab.com/build-graph/connectiors"
	"gitlab.com/build-graph/models"
	pb "gitlab.com/restra-core/buffers/vendors"
)

var cfg *models.Config

func init() {
	cfg = &models.Config{
		//Couchbase information
		DBHost:         "139.180.219.37",
		DBPort:         "8091",
		NoSqlUser:      "Administrator",
		NoSqlPassword:  "restra247",
		BucketName:     "restra2",
		BucketPassword: "",

		//Bolt information
		BoltHost:     "localhost",
		BoltPort:     "7687",
		BoltUser:     "neo4j",
		BoltPassword: "n4j",
	}
}

func main() {

	couchbase, err := connectiors.GetDBConnection(cfg)
	if err != nil {
		panic(err)
	}

	//bolt , err := connectiors.GetBoltConnection(cfg)
	//if err != nil {
	//	panic(err)
	//}

	query := gocb.NewN1qlQuery(`SELECT * FROM restra2 WHERE type="restaurant"`)
	rows, err := couchbase.ExecuteN1qlQuery(query, nil)

	var row map[string]interface{}
	for rows.Next(&row) {
		//Serializing the original restaurant json
		b, err := json.Marshal(row["restra2"])
		if err != nil {
			panic(err)
		}

		//Just for Printing out
		var restaurant pb.RestaurantInfo
		err = json.Unmarshal(b, &restaurant)
		if err != nil {
			panic(err)
		}
		fmt.Println(restaurant)

		BuildGraph(b, &restaurant)
	}
	if err = rows.Close(); err != nil {
		fmt.Printf("Couldn't get all the rows: %s\n", err)
	}
}

func BuildGraph(js []byte, restaurant *pb.RestaurantInfo) {
	var data map[string]interface{}

	err := json.Unmarshal(js, &data)
	if err != nil {
		panic(err.Error())
	}

	j := J.Jypher{
		ParentNode: &jmodel.EntityInfo{
			Name: "restaurant",
			ID:   restaurant.Id,
		},
	}
	jsonInfo := jmodel.JSONInfo{
		DecodedJSON: data,
		Rules: &jmodel.Rules{
			SkipField: []string{"picture", "description", "updated_at", "created_at"},
		},
	}
	decodedGraph, err := j.GetJypher(jsonInfo)
	if err != nil {
		panic(err.Error())
	}

	cypher := j.BuildCypher(decodedGraph)

	bolt, _ := connectiors.GetBoltConnection(cfg)

	for i, v := range cypher {
		if i >= 10{
			break
		}
		// Start by creating a node
		fmt.Println(v)
		_, err := bolt.Connection.ExecNeo(v, nil)
		if err != nil {
			panic(err.Error())
		}
	}
}
