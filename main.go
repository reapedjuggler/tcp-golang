package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
)

type HttpMethod string
type StatusCode string

type HandlerFunc func(request *http.Request) (StatusCode, string)

type Endpoints struct {
	endpointsMap map[string]map[HttpMethod]HandlerFunc
}

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PATCH  HttpMethod = "PATCH"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

const (
	OK               StatusCode = "200 OK"
	MethodNotAllowed StatusCode = "405 Method Not Allowed"
	NotFound         StatusCode = "404 Not Found"
)

func (endpoints *Endpoints) AddRoute(route string, method HttpMethod, handler HandlerFunc) {

}

// listen on a port using net or sth
// accept the connection on a port for an infinite time
// accept the connection
// create a goroutine for each connection
// handle the connection
// try removing the go keyword and see, ideally it should work in a sync manner

func main() {
	listener, err := net.Listen("tcp", ":8000")
	defer listener.Close()
	log.Println("Server listening on :8000")

	if err != nil {
		log.Print("Failed to listen at the port", err)
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Failed to accept the connection", err)
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// read the request from the connection
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Print("Failed to read the request", err)
		panic(err)
	}
	fmt.Printf("Request %s %s\n", request.Method, request.URL)
	// write the response to the connection
	responseBody := "Lets go, " + request.URL.Path + "!\r\n"
	headers := fmt.Sprintf("HTTP/1.1 %s\r\n", OK)
	headers += fmt.Sprintf("Content-Length: %d\r\n", len(responseBody))
	headers += "Content-Type: text/plain\r\n\r\n"

	// write the headers and body to the connection
	conn.Write([]byte(headers))
	conn.Write([]byte(responseBody))

	fmt.Print("Response sent\n")
}
