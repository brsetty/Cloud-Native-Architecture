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
	mongodbEndpoint = "mongodb://192.168.85.128:32619" // 172.18.0.2
)

type Item struct {
	Name  string  `bson:"name"`
	Price float64 `bson:"price"`
}

func main() {

	http.HandleFunc("/list", list)
	http.HandleFunc("/price", price)
	http.HandleFunc("/add", add)
	http.HandleFunc("/update", update)
	http.HandleFunc("/delete", delete)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// function creates an entry
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

	//creates a database shopDB and collection items in mongodb
	itemscollection := client.Database("shopDB").Collection("items")
	//filters by name. upon giving db.items.find() everthing is listed.
	filter := bson.M{"name": bson.M{"$eq": itm}}

	result := itemscollection.FindOne(context.TODO(), filter).Decode(&Item{})
	if result == nil {
		fmt.Fprintf(w, " %s already exists in database\n", itm)
	} else {
		//if it is a new entry, InsertOne is used. It is an api for communicating with mongo
		_, err = itemscollection.InsertOne(context.TODO(), &Item{
			Name:  itm,
			Price: price,
		})
		checkError(err)
		fmt.Fprintf(w, "Item %s added sucessfully\n", itm)
	}
}

// list all the items in the collection of our database
func list(w http.ResponseWriter, req *http.Request) {
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shopDB").Collection("items")
	result, err := itemscollection.Find(context.TODO(), bson.D{{}})
	checkError(err)
	fmt.Fprintf(w, "Listing items in database ....\n")
	for result.Next(context.TODO()) {
		var res Item
		err = result.Decode(&res)
		checkError(err)
		fmt.Fprintf(w, "%s: %f\n", res.Name, res.Price)
	}
}

func price(w http.ResponseWriter, req *http.Request) {
	itm := req.URL.Query().Get("item")
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shopDB").Collection("items")
	filter := bson.M{"name": bson.M{"$eq": itm}}
	var i *Item
	err = itemscollection.FindOne(context.TODO(), filter).Decode(i)
	if err != nil {
		fmt.Fprintf(w, "%s Item not found", i.Name)
	} else {
		fmt.Fprintf(w, "%f\n", i.Price)
	}

}

// Updates the pre existing entry in database
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

	itemscollection := client.Database("shopDB").Collection("items")

	result, err := itemscollection.UpdateOne(
		context.TODO(),
		bson.M{"name": itm},
		bson.D{
			{"$set", bson.D{{"price", price}}},
		},
	)
	checkError(err)
	fmt.Fprintf(w, "Updated price of %d of item %s\n", result.ModifiedCount, itm)
}

// Deletes the entry from database
func delete(w http.ResponseWriter, req *http.Request) {

	// get item from the URL query parametes
	params := req.URL.Query()
	itm := params.Get("item")
	client := getClient()
	err := client.Connect(context.TODO())
	checkError(err)
	defer client.Disconnect(context.TODO())
	itemscollection := client.Database("shopDB").Collection("items")

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
