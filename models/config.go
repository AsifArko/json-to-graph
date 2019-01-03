package models

type Config struct {
	//Serve conf
	ServePort string

	//Couchbase conf
	DBHost         string
	DBPort         string
	NoSqlUser      string
	NoSqlPassword  string
	BucketName     string
	BucketPassword string

	//Bolt conf
	BoltHost     string
	BoltPort     string
	BoltUser     string
	BoltPassword string

	//Redis conf
	RedisHost     string
	RedisPort     string
	RedisUser     string
	RedisPassword string
}
