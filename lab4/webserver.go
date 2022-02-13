package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type database struct {
	data map[string]dollars
	sync.RWMutex
}

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 5}}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/delete", db.delete)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

//type database map[string]dollars

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "List of items available\n")
	db.RLock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	db.RUnlock()
}
func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db.data[item]; ok {
		fmt.Fprintf(w, "The price of %s is : %s\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	//db.Lock()
	if _, ok := db.data[item]; ok {
		fmt.Fprintf(w, "%s is deleted\n", item)
		db.Lock()
		delete(db.data, item)
		db.Unlock()
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	//db.Unlock()
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	db.Lock()
	if _, ok := db.data[item]; ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "You have an item already: %q\n", item)
	} else {
		price1, e := strconv.ParseFloat(price, 32)
		if e != nil {
			fmt.Fprintf(w, "You have an invalid price\n")
		} else {
			db.data[item] = dollars(price1)
			fmt.Fprintf(w, "%s is added to the list\n", item)
		}
	}
	db.Unlock()
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	db.Lock()
	if _, ok := db.data[item]; ok {
		price1, e := strconv.ParseFloat(price, 32)
		if e != nil {
			fmt.Fprintf(w, "You have an invalid price\n")
		} else {
			db.data[item] = dollars(price1)
			fmt.Fprintf(w, "The price of the %s is updated to %s\n", item, price)
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
	db.Unlock()
}
