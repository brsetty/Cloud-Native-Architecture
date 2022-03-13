// Package main implements a server for movieinfo service.
package main

import (
	"LABS/lab_5/movieapi"
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
type database struct {
	moviedb map[string][]string
	sync.RWMutex
}

var db database

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	title := in.GetTitle()
	log.Printf("Received GET: %v", title)
	reply := &movieapi.MovieReply{}

	db.Lock()
	defer db.Unlock()

	if val, ok := db.moviedb[title]; !ok { // Title not present in database
		return reply, errors.New("Not Present in Database")
	} else {
		if year, err := strconv.Atoi(val[0]); err != nil {
			reply.Year = -1
		} else {
			reply.Year = int32(year)
		}
		reply.Director = val[1]
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)

	}

	return reply, nil

}

func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.Status, error) {
	var castName string
	status := &movieapi.Status{}

	title := in.GetTitle()
	log.Printf("Received SET: %v", title)

	year := in.GetYear()
	director := in.GetDirector()
	cast := in.GetCast()

	if (title == "") || (director == "") || (year <= 0) || (len(cast) == 0) { //Checking for Invalid Input
		status.Code = "SET" + " " + title + " " + "FAIL!!!"
		return status, errors.New("ERROR: Invalid Input!!!")
	} else { //Valid Input Provided
		db.Lock()
		defer db.Unlock()

		if _, ok := db.moviedb[title]; ok {
			status.Code = "SET" + " " + title + " " + "FAIL!!!"
			return status, errors.New("ERROR: Title is already present in database!!!")
		} else {
			db.moviedb[title] = append(db.moviedb[title], strconv.FormatInt(int64(year), 10))
			db.moviedb[title] = append(db.moviedb[title], director)

			for i, name := range cast {
				if i == 0 {
					castName = name + ","
				} else if i == len(cast)-1 {
					castName = castName + name
				} else {
					castName = castName + name + ","
				}
			}

			db.moviedb[title] = append(db.moviedb[title], castName)

			status.Code = "SET" + " " + title + " " + "SUSCCESS!!!"

			return status, nil
		}
	}
}

func main() {
	db.Lock()
	defer db.Unlock()
	db = database{moviedb: map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
