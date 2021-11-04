package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"time"
)

type MongoConf struct {
	MongoRunMode string

	MongoUriSingle       string
	MongoUriReplicaSet   string
	MongoDBName          string
	MongoUser            string
	MongoPwd             string
	MongoMaxPoolSize     uint64
	MongoMaxConnIdleTime int
}

func main() {
	ctx := context.Background()

	getMongoCliByReplicaSet(ctx, &MongoConf{MongoUser: "root", MongoPwd: "Ft_Mongo_123",
		MongoUriReplicaSet: "mongodb://172.16.71.17:27017,172.16.71.17:27018,172.16.71.17:27019/?replicaSet=myset"})

}
func getMongoCliByReplicaSet(ctx context.Context, mongoConfObj *MongoConf) *mongo.Client {
	// Set the URI of ReplicaSet MongoDB to connect to Mongodb
	fmt.Println(mongoConfObj.MongoUriReplicaSet)
	opt := options.Client().ApplyURI(mongoConfObj.MongoUriReplicaSet)
	// Set user and password
	fmt.Println(mongoConfObj.MongoUser, mongoConfObj.MongoPwd)
	credential := options.Credential{
		Username: mongoConfObj.MongoUser,
		Password: mongoConfObj.MongoPwd,
	}
	opt.SetAuth(credential)
	// Set the maximum number of connections using the connection pool
	opt.SetMaxPoolSize(mongoConfObj.MongoMaxPoolSize)
	// Set the maximum time the connection can remain idle
	opt.SetMaxConnIdleTime(time.Duration(mongoConfObj.MongoMaxConnIdleTime) * time.Second)
	// Set the read preference to use for read operations
	wantReadPref, err := readpref.New(readpref.SecondaryPreferredMode)
	if err != nil {
		log.Panic(err)
	}
	opt.SetReadPreference(wantReadPref)
	// Specifies that the query should return the instance’s most recent data
	// acknowledged as having been written to a majority of members in the replica set.
	opt.SetReadConcern(readconcern.Majority())
	// WMajority requests acknowledgement that write operations
	// propagate to the majority of mongodb instances.
	wantWriteConcern := writeconcern.New(writeconcern.WMajority())
	opt.SetWriteConcern(wantWriteConcern)

	cli, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Panic("Failed to connect ReplicaSet MongoDB, error message: ", err.Error())
		return nil
	}

	err = cli.Ping(ctx, nil)
	if err != nil {
		log.Panic("Failed to execute MongoDB client ping command, error message: ", err.Error())
		return nil
	}

	_, err = cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Panic("Test MongoDB authentication failed, error message:", err.Error())
		return nil
	}

	log.Println("Successfully to get Mongodb client of ReplicaSet")
	return cli
}
