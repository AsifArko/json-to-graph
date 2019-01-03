package connectiors

import (
	"errors"
	"fmt"
	"github.com/couchbase/gocb"
	"gitlab.com/build-graph/models"
)

type DBConnection struct {
	*gocb.Bucket
}

func GetDBConnection(cfg *models.Config) (*DBConnection, error) {

	cluster, err := gocb.Connect("couchbase://" + cfg.DBHost)
	if err != nil {
		msg := fmt.Sprintf("DBConnection :: Couchbase :: Cluster :: Error : %s", err.Error())
		fmt.Println(msg)
		return nil, errors.New(msg)
	}

	auth := &gocb.PasswordAuthenticator{
		Username: cfg.NoSqlUser,
		Password: cfg.NoSqlPassword,
	}

	err = cluster.Authenticate(auth)

	if err != nil {
		msg := fmt.Sprintf("DBConnection :: Couchbase :: Cluster > Auth :: Error : %s", err.Error())
		fmt.Println(msg)
		return nil, errors.New(msg)
	}

	var bucket *gocb.Bucket

	// It takes some time to connect to couchbase bucket thats why
	for {
		bucket, err = cluster.OpenBucket(cfg.BucketName, cfg.BucketPassword)
		if err == nil {
			break
		}
		//msg := fmt.Sprintf("DBConnection :: Couchbase :: Cluster > Open Bucket :: Error : %s :: trying again ...", err.Error())
		//fmt.Println(msg)
	}

	fmt.Println("\n\nConnected to NOSQL DB")

	return &DBConnection{bucket}, nil
}
