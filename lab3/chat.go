// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	name    string
	channel chan<- string // an outgoing message channel
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			active_client(clients, cli.channel)
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.channel)
		}
	}
}

func handleConn(conn net.Conn) {

	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	who := conn.RemoteAddr().String()
	ch <- "enter your name"
	name := input.scan
	ch <- "You are " + who
	messages <- who + " has arrived"
	cl := client{who, ch}
	entering <- cl

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cl
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

func active_client(clients map[client]bool, channel chan<- string) {
	if len(clients) == 0 {
		channel <- "There are no other users in the network"
		return
	}
	channel <- "All the Active Users are: "
	for cli := range clients {
		channel <- cli.name
	}
}
