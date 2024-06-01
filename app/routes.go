package main

import (
	"net"
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
		response := formatResponse(Response{
			status: 404,
			reason: "Not Found",
			header: map[string]string{},
			body:   "",
		})
		c.Write([]byte(response))
	}
}
