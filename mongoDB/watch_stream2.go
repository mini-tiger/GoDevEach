package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	//"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

var uri string = "mongodb://cc:cc@172.22.50.25:27117,172.22.50.25:27118/?authMechanism=SCRAM-SHA-256&authSource=cmdb&directConnection=true"

func main() {
	// Set up client options
	clientOpts := options.Client().ApplyURI(uri)
	// Set up the client with write concern
	//clientOpts.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(err)
	}

	// Ping the primary to ensure you have a connection
	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		panic(err)
	}

	// Set up the database and collection names
	dbName := "cmdb"
	collName := "cc_ObjectBase"

	// Set up the pipeline for the Change Stream
	pipeline := mongo.Pipeline{}

	// Set up the options for the Change Stream
	options := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	options.SetMaxAwaitTime(time.Second * 5) // xxx 指定新文档游标上一个的的最长等待时间，防止 重复watch到数据
	// Set up the resume token to start monitoring from the current time
	//resumeToken := &bson.Raw{bson.M{"_id": "resume_token_id"}}
	fmt.Println(options)
	// Start the Change Stream
	changeStream, err := client.Database(dbName).Collection(collName).Watch(
		context.Background(), pipeline, options)

	if err != nil {
		panic(err)
	}

	// Set up a loop to print out the Change Stream events
	for changeStream.Next(context.Background()) {
		var changeDoc bson.M
		if err := changeStream.Decode(&changeDoc); err != nil {
			panic(err)
		}

		fmt.Println("Received change event:", changeDoc)
		// Set the resume token for the next Change Stream event
		//*resumeToken = changeStream.ResumeToken()

		// Wait for a second before printing the next Change Stream event
		time.Sleep(time.Second)
	}

	// Check for any errors that may have occurred while processing the Change Stream
	if err := changeStream.Err(); err != nil {
		panic(err)
	}
}
