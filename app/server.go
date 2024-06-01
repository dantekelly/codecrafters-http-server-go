package main

import (
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	c, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Connection accepted")
	response := formatResponse(Status{200, "OK"})
	c.Write([]byte(response))
}

type Status struct {
	status int
	reason string
}

func formatResponse(status Status) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", status.status, status.reason)
	// fmt.Printf("Content-Length: %d\n", len(body))
	// fmt.Printf("Content-Type: text/plain\n\n")
	// fmt.Println(body)
}
