package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 5}}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/add", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	data map[string]dollars
	sync.RWMutex
}

// create
func (db *database) add(w http.ResponseWriter, req *http.Request) {
	// get the item and price from the URL query parameters
	params := req.URL.Query()

	item := params.Get("item")
	price, err := strconv.ParseFloat(params.Get("price"), 32)
	// check if price is valid
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "invalid price: %f\n", price)
		log.Print(err)
	}

	// aquire the lock and add it to the database
	db.Lock()
	db.data[item] = dollars(price)
	db.Unlock()

	fmt.Fprintf(w, "Item %s : %s added to the database\n", item, dollars(price))
}

// Read
func (db *database) list(w http.ResponseWriter, req *http.Request) {
	db.RLock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	db.RUnlock()
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	db.RLock()
	price, ok := db.data[item]
	db.RUnlock()

	if ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

// Update
func (db *database) update(w http.ResponseWriter, req *http.Request) {
	// get item and price from the URL query parameters
	params := req.URL.Query()
	item := params.Get("item")
	price, err := strconv.ParseFloat(params.Get("price"), 32)

	// check if price is valid
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "invalid price: %f\n", price)
		log.Print(err)
	}
	// aquire the db lock
	db.Lock()
	// check if the item is present in the database if yes update the price
	if _, ok := db.data[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "Item %s Not Found\n", item)
		log.Printf("item %s not found in the database", item)
	} else {
		db.data[item] = dollars(price)
		fmt.Fprintf(w, "Item %s : %s updated to the database\n", item, dollars(price))
	}
	// release the db lock
	db.Unlock()

}

// Delete
func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	// get item from the URL query parametes
	params := req.URL.Query()
	item := params.Get("item")

	// aquire the lock
	db.Lock()
	// check if the item is in the database
	if _, ok := db.data[item]; !ok {
		w.WriteHeader(http.StatusBadRequest) // 404
		fmt.Fprintf(w, "Item %s Not Found\n", item)
		log.Printf("item %s not found in the database", item)
	} else {
		delete(db.data, item)
		fmt.Fprintf(w, "Item %s deleted from the database\n", item)
	}
	// release the lock
	db.Unlock()

}
