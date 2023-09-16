package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"noelzubin/redis-go/protocol"
)

func main() {
	// Connect to Redis
	host := getHostName()
	conn, err := net.Dial("tcp", host)

	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	}
	defer conn.Close()

	// Test the connection
	conn.Write(encodeString("PING"))
	resp, err := protocol.DecodeRESP(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error testing Redis connection:", err)
		return
	}

	if !strings.HasPrefix(resp.Output(), "PONG") {
		fmt.Println("Unexpected Redis response:", resp)
		return
	}

	// Handling Ctrl-C gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nReceived Ctrl-C. Exiting...")
		conn.Close()
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Print the prompt
		fmt.Print("> ")

		// Read user input
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Skip empty lines
		if strings.TrimSpace(input) == "" {
			continue
		}

		// Send commands to redis
		conn.Write(encodeString(strings.TrimSpace(input)))

		// Get response
		buf, err := protocol.DecodeRESP(bufio.NewReader(conn))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Println(buf.Output())
	}
}

func encodeString(s string) []byte {
	r := strings.Split(s, " ")
	sa := protocol.NewArrayStringValue(r)
	return sa.Encode()
}

func getHostName() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	return "localhost:6379"
}
