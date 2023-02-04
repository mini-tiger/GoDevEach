package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// InstAsst an association definition between instances.
type InstAsst struct {
	// sequence ID
	ID int64 `field:"id" json:"id,omitempty"`
	// inst id associate to ObjectID
	InstID int64 `field:"bk_inst_id" json:"bk_inst_id,omitempty" bson:"bk_inst_id"`
	// association source ObjectID
	ObjectID string `field:"bk_obj_id" json:"bk_obj_id,omitempty" bson:"bk_obj_id"`
	// inst id associate to AsstObjectID
	AsstInstID int64 `field:"bk_asst_inst_id" json:"bk_asst_inst_id,omitempty"  bson:"bk_asst_inst_id"`
	// association target ObjectID
	AsstObjectID string `field:"bk_asst_obj_id" json:"bk_asst_obj_id,omitempty" bson:"bk_asst_obj_id"`
	// bk_supplier_account
	OwnerID string `field:"bk_supplier_account" json:"bk_supplier_account,omitempty" bson:"bk_supplier_account"`
	// association id between two object
	ObjectAsstID string `field:"bk_obj_asst_id" json:"bk_obj_asst_id,omitempty" bson:"bk_obj_asst_id"`
	// association kind id
	AssociationKindID string `field:"bk_asst_id" json:"bk_asst_id,omitempty" bson:"bk_asst_id"`

	// BizID the business ID
	BizID int64 `field:"bk_biz_id" json:"bk_biz_id,omitempty" bson:"bk_biz_id"`
}

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

// This is an example change event struct for inserts.
// It does not include all possible change event fields.
// You should consult the change event documentation for more info:
// https://docs.mongodb.com/manual/reference/change-events/
type changeEvent struct {
	ID            changeID            `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  InstAsst            `bson:"fullDocument"`
	DocumentKey   documentKey         `bson:"documentKey"`
	Ns            namespace           `bson:"ns"`
}

//var uri string = "mongodb://cc:cc@172.22.50.25:27021/?connectTimeoutMS=300000&authSource=cmdb"
var uri string = "mongodb://cc:cc@172.22.50.25:32047/?connectTimeoutMS=300000&authSource=cmdb"

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("cmdb").Collection("cc_InstAsst")

	// Watches the todo collection and prints out any changed documents
	go watch(collection)
	select {}
	// Inserts random todo items at two second intervals
	//insert(collection)

}

func watch(collection *mongo.Collection) {

	// xxx FullDocument 返回由更新操作修改的那些字段的增量，而不是整个更新的文档。 https://www.mongodb.com/docs/manual/reference/method/db.collection.watch/
	streamOptions := options.ChangeStream().SetFullDocument(options.WhenAvailable)
	//streamOptions := options.ChangeStream().SetFullDocumentBeforeChange(options.WhenAvailable) //mongodb 6.x

	/* xxx fullDocument删除文档时不返回源文档
	https://stackoverflow.com/questions/56939610/how-to-get-fulldocument-from-mongodb-changestream-when-a-document-is-deleted
	*/

	// Watch the  collection
	cs, err := collection.Watch(context.TODO(), mongo.Pipeline{}, streamOptions)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent interface{}

		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Change Event: %+v\n", changeEvent)

	}

}
