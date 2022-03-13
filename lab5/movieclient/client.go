// Package main imlements a client for movieinfo service
package main

import (
	"LABS/lab_5/movieapi"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var (
	title    string
	year     int32
	director string
	cast     []string
)

func set() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	status, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: title, Year: year, Director: director, Cast: cast})
	if err != nil {
		errMsg := "SET" + " " + title + ":"
		log.Fatalf("%s %v", errMsg, err)
	}
	log.Printf("%s", status.GetCode())
}

func get() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	t := title

	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: t})
	if err != nil {
		errMsg := "GET" + " " + title + ":"
		log.Fatalf("%s %v", errMsg, err)
	}
	log.Printf("Movie Info for %s %d %s %v", t, r.GetYear(), r.GetDirector(), r.GetCast())
}

func main() {
	var service string
	var name []string
	var yr string

	service = os.Args[1]

	if service == "1" {
		name = strings.Split(os.Args[2], " ")
		for i, word := range name {
			if i == 0 {
				title = word
			} else if i == len(cast)-1 {
				title = title + " " + word
			} else {
				title = title + " " + word
			}
		}
		get()

	} else {
		name = strings.Split(os.Args[2], "_")
		for i, word := range name {
			if i == 0 {
				title = word
			} else if i == len(cast)-1 {
				title = title + " " + word
			} else {
				title = title + " " + word
			}
		}

		yr = os.Args[3]
		i, _ := strconv.Atoi(yr)
		year = int32(i)

		name = strings.Split(os.Args[4], "_")
		for i, word := range name {
			if i == 0 {
				director = word
			} else if i == len(cast)-1 {
				director = director + " " + word
			} else {
				director = director + " " + word
			}
		}

		cast = strings.Split(os.Args[5], ",")

		set()
	}
}
