package main

import (
	"LABS/lab_5/movieapi"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

const (
	address = ":8050"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	Input := bufio.NewScanner(os.Stdin)
	fmt.Println("To add a movie: 1 , To search a movie : 2 ")
	Input.Scan()
	opt := Input.Text()
	option, _ := strconv.ParseInt(opt, 10, 32)
	if option == 1 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("enter the movie name,year,director,cast")
		scanner.Scan()
		text := scanner.Text()
		names := strings.Split(text, ",")
		year, _ := strconv.ParseInt(names[1], 10, 32)
		newMovie := &movieapi.MovieData{
			Title:    names[0],
			Year:     int32(year),
			Director: names[2],
			Cast:     names[3:],
		}

		addMovie(c, newMovie)
	} else {
		Movietitle := bufio.NewScanner(os.Stdin)
		fmt.Println("Movie: ")
		Movietitle.Scan()
		title := Movietitle.Text()
		getMovie(c, title)
	}

}

func addMovie(client movieapi.MovieInfoClient, movie *movieapi.MovieData) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SetMovieInfo(ctx, movie)
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie %s :%s", movie.Title, resp.Code)

}

func getMovie(client movieapi.MovieInfoClient, title string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s", title)
	log.Printf("Year:%d", resp.GetYear())
	log.Printf("Movie Director:%s", resp.GetDirector())
	log.Printf("Movie Cast:%v", resp.GetCast())
}
