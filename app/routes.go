package main

import (
	"log"
	"net"
	"os"
	"strings"
)

func handleRoutes(c net.Conn, request Request) {
	switch {
	case strings.Contains(request.path, "/echo/"):
		route := strings.Split(request.path, "/echo/")[1]
		response := formatResponse(Response{
			status: 200,
			reason: "OK",
			header: map[string]string{
				"Content-Type": "text/plain",
			},
			body: route,
		}, request.header.acceptEncoding)
		c.Write([]byte(response))
	case request.path == "/user-agent":
		response := formatResponse(Response{
			status: 200,
			reason: "OK",
			header: map[string]string{
				"Content-Type": "text/plain",
			},
			body: request.header.agent,
		}, "")
		c.Write([]byte(response))
	case request.path == "/":
		response := formatResponse(Response{
			status: 200,
			reason: "OK",
			header: map[string]string{},
			body:   "",
		}, "")
		c.Write([]byte(response))
	default:
		if request.method == "POST" {
			postFileRoute(c, request)
			return
		}

		getFileRoute(c, request.path)
	}
}

func getFileRoute(c net.Conn, path string) {
	path = strings.TrimPrefix(path, "/files/")

	if config.directory != "" {
		path = config.directory + path
	}

	file, err := os.ReadFile(path)
	if err != nil {
		response := formatResponse(Response{
			status: 404,
			reason: "Not Found",
			header: map[string]string{},
			body:   "",
		}, "")
		c.Write([]byte(response))
		return
	}

	r := formatResponse(Response{
		status: 200,
		reason: "OK",
		header: map[string]string{
			"Content-Type": "application/octet-stream",
		},
		body: string(file),
	}, "")

	c.Write([]byte(r))
}

func postFileRoute(c net.Conn, request Request) {
	path := strings.TrimPrefix(request.path, "/files/")

	if config.directory != "" {
		path = config.directory + path
	}

	file := []byte(strings.Trim(request.body, "\x00"))

	err := os.WriteFile(path, file, 0644)
	if err != nil {
		log.Printf("Error writing file: %v\n", err)

		response := formatResponse(Response{
			status: 500,
			reason: "Internal Server Error",
			header: map[string]string{},
			body:   "",
		}, "")
		c.Write([]byte(response))
		return
	}

	response := formatResponse(Response{
		status: 201,
		reason: "Created",
		header: map[string]string{},
		body:   "",
	}, "")
	c.Write([]byte(response))
}
