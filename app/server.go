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

var config Config

func init() {
	directory := ""
	// Gather arguments
	args := os.Args[1:]
	if len(args) > 0 {
		for i, arg := range args {
			if arg == "--directory" {
				if i+1 < len(args) {
					directory = args[i+1]
				} else {
					log.Println("Directory flag provided without a directory")
				}
			}
		}
	}

	builder := ConfigBuilder{}
	builder.Directory(directory)
	cfg, err := builder.Build()
	if err != nil {
		log.Println("Failed to build config: ", err.Error())
		os.Exit(1)
	}

	config = *cfg
}

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
	headers := rawRequest[1 : len(rawRequest)-2]
	body := rawRequest[len(rawRequest)-1]
	headersMap := make(map[string]string)

	for _, header := range headers {
		headerSlice := strings.Split(header, ": ")

		headersMap[headerSlice[0]] = headerSlice[1]
	}

	request := Request{
		method:  requestSlice[0],
		path:    requestSlice[1],
		version: requestSlice[2],
		header: RequestHeader{
			host:           headersMap["Host"],
			agent:          headersMap["User-Agent"],
			accept:         headersMap["Accept"],
			acceptEncoding: headersMap["Accept-Encoding"],
		},
		body: body,
	}

	// log.Printf("Request: %v\n", body)

	handleRoutes(c, request)
}

type RequestHeader struct {
	host           string
	agent          string
	accept         string
	acceptEncoding string
}
type Request struct {
	method  string
	path    string
	version string
	header  RequestHeader
	body    string
}

func formatRequest(r Request) string {
	return fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nAccept: %s\r\n\r\n", r.method, r.path, r.header.host, r.header.agent, r.header.accept)
}

type Response struct {
	status int
	reason string
	header map[string]string
	body   string
}

func formatResponse(r Response, e string) string {
	// Start with the status line
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.status, r.reason)

	// Add headers
	for key, value := range r.header {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	if r.body != "" {
		response += fmt.Sprintf("Content-Length: %d\r\n", len(r.body))
	}
	if e == "gzip" {
		response += fmt.Sprintf("Content-Encoding: %s\r\n", e)
	}

	// Add a blank line to indicate the end of headers
	response += "\r\n"

	if r.body != "" {
		response += r.body
	}

	return response
}
