package main

import (
	"fmt"
	"net"
	"noelzubin/redis-go/eventloop"
	"noelzubin/redis-go/store"
	"os"
)

func main() {
	store := store.InitStore()
	el := eventloop.InitEventloop(store)

	// Main Event loop
	go el.RunLoop()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	fmt.Println("running server on port 6379")

	for {
		conn, err := l.Accept()
		fmt.Println("accepted connection ")

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("accepted connection ")
		go el.HandleConnection(conn)
	}
}
