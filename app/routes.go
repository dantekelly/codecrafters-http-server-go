package main

import (
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
		})
		c.Write([]byte(response))
	case request.path == "/user-agent":
		response := formatResponse(Response{
			status: 200,
			reason: "OK",
			header: map[string]string{
				"Content-Type": "text/plain",
			},
			body: request.header.agent,
		})
		c.Write([]byte(response))
	case request.path == "/":
		response := formatResponse(Response{
			status: 200,
			reason: "OK",
			header: map[string]string{},
			body:   "",
		})
		c.Write([]byte(response))
	default:
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
		})
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
	})

	c.Write([]byte(r))
}
