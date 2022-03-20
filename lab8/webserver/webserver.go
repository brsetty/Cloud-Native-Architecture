package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://0.0.0.0:27017" // 172.18.0.2
)

type Item struct {
	Item  string  `bson:"item"`
	Price float64 `bson:"price"`
}

func main() {

	http.HandleFunc("/list", list)
	http.HandleFunc("/price", price)
	http.HandleFunc("/add", add)
	//http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	log.Fatal(http.ListenAndServe(":8050", nil))
}

// create
func add(w http.ResponseWriter, req *http.Request) {
	// get the Item and price from the URL query parameters
	params := req.URL.Query()

	itm := params.Get("item")
	price, err := strconv.ParseFloat(params.Get("price"), 32)
	checkError(err)
	client := getClient()
	err = client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())

	itemscollection := client.Database("shop").Collection("stock")
	filter := bson.M{"item": bson.M{"$eq": itm}}
	result := itemscollection.FindOne(context.TODO(), filter).Decode(&Item{})
	if result == nil {
		fmt.Fprintf(w, "Already Exists\n")
	} else {
		_, err = itemscollection.InsertOne(context.TODO(), &Item{
			Item:  itm,
			Price: price,
		})
		checkError(err)
		fmt.Fprintf(w, "Item %s added\n", itm)
	}
}

// Read
func list(w http.ResponseWriter, req *http.Request) {
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shop").Collection("stock")
	result, err := itemscollection.Find(context.TODO(), bson.D{{}})
	checkError(err)
	for result.Next(context.TODO()) {
		var res Item
		err = result.Decode(&res)
		checkError(err)
		fmt.Fprintf(w, "%s: %f\n", res.Item, res.Price)
	}
}

func price(w http.ResponseWriter, req *http.Request) {
	itm := req.URL.Query().Get("item")
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shop").Collection("stock")
	filter := bson.M{"item": bson.M{"$eq": itm}}
	var i *Item
	err = itemscollection.FindOne(context.TODO(), filter).Decode(i)
	if err != nil {
		fmt.Fprintf(w, "Item Not found")
	} else {
		fmt.Fprintf(w, "%f\n", i.Price)
	}

}

/* Update
func update(w http.ResponseWriter, req *http.Request) {
	// get Item and price from the URL query parameters
	params := req.URL.Query()
	itm := params.Get("item")
	price, err := strconv.ParseFloat(params.Get("price"), 64)

	// check if price is valid
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "invalid price: %f\n", price)
		log.Print(err)
	}
	client := getClient()
	err = client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())

	itemscollection := client.Database("shop").Collection("stock")

	result, err := itemscollection.UpdateOne(
		context.TODO(),
		bson.M"item": itm},
		bson.D{
			{"$set", bson.D{"price", price}},
		},
	)
	checkError(err)
	fmt.Fprintf(w, "Updated Price of %d item\n", result.ModifiedCount)
}*/

// Delete
func delete(w http.ResponseWriter, req *http.Request) {
	// get item from the URL query parametes
	params := req.URL.Query()
	itm := params.Get("item")
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shop").Collection("stock")

	result, err := itemscollection.DeleteOne(context.TODO(), bson.M{"name": itm})
	checkError(err)

	fmt.Printf("DeleteOne removed %v document(s)\n", result.DeletedCount)
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getClient() *mongo.Client {
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)
	return client
}
