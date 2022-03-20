// Example use of Go mongo-driver
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://127.0.0.1:27017" // Find this from the Mongo container
)

/*type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string             `bson:"title"`
	Body      string             `bson:"body"`
	Tags      []string           `bson:"tags"`
	Comments  uint64             `bson:"comments"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}*/
type Post struct {
	Item  string  `bson:"item"`
	Price float64 `bson:"price"`
}

func main() {
	// create a mongo client
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	// Connect to mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	checkError(err)
	// Disconnect
	defer client.Disconnect(ctx)

	// select collection from database
	col := client.Database("shop").Collection("stock")
	filter := bson.M{"item": bson.M{"$elemMatch": bson.M{"$eq": "shoes"}}}
	var p Post
	result := col.FindOne(ctx, filter).Decode(&p)
	fmt.Printf("%v \n", result)
	// Insert one
	res, err := col.InsertOne(ctx, &Post{
		Item:  "shoes",
		Price: 100.50,
	})
	checkError(err)
	fmt.Printf("inserted item: %s\n", res.InsertedID)

	// filter posts tagged as mongod
	// filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "mongodb"}}}

	// find one document
	if err = col.FindOne(ctx, filter).Decode(&p); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("post: %+v\n", p)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
