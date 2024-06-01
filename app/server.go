package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	log.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		log.Println("Connection accepted")
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		log.Println("Error reading:", err.Error())
		return
	}

	rawRequest := strings.Split(string(buf), "\r\n")

	requestSlice := strings.Split(rawRequest[0], " ")
	request := Request{
		method:  requestSlice[0],
		path:    requestSlice[1],
		version: requestSlice[2],
	}
	log.Printf("Request: %+v", request)

	//headersSlice := strings.Split(rawRequest[1], " ")
	// request := Request{}
	// log.Printf("Headers: %+v", headersSlice)

	switch request.path {
	case "/":
		response := formatResponse(Response{200, "OK"})
		c.Write([]byte(response))
	default:
		response := formatResponse(Response{404, "Not Found"})
		c.Write([]byte(response))
	}

}

type RequestHeader struct {
	host   string
	agent  string
	accept string
}
type Request struct {
	method  string
	path    string
	version string
	header  RequestHeader
}

func formatRequest(r Request) string {
	return fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nAccept: %s\r\n\r\n", r.method, r.path, r.header.host, r.header.agent, r.header.accept)
}

type Response struct {
	status int
	reason string
}

func formatResponse(r Response) string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", r.status, r.reason)
	// fmt.Printf("Content-Length: %d\n", len(body))
	// fmt.Printf("Content-Type: text/plain\n\n")
	// fmt.Println(body)
}
