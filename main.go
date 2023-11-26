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

var endpoints *Endpoints

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
	val := endpoints.endpointsMap[route]
	log.Println("Inside AddRoute")
	if len(val) == 0 {
		fmt.Println("No definition exists for this route, hence create a new one")
		endpoints.endpointsMap[route] = map[HttpMethod]HandlerFunc{}
	}
	// Check if this gives an error
	endpoints.endpointsMap[route][method] = handler
}

func (endpoints *Endpoints) registerEndpoint() {

	log.Println("Inside registerEndpoint")
	endpoints.AddRoute("/", GET, func(request *http.Request) (StatusCode, string) {
		body := "Bonjour VT, Route 0 here!"
		return OK, body
	})

	endpoints.AddRoute("/route-1", GET, func(request *http.Request) (StatusCode, string) {
		body := "Bonjour VT, Route 1 here!"
		return OK, body
	})

	endpoints.AddRoute("/post", POST, func(request *http.Request) (StatusCode, string) {
		buf := make([]byte, 1024)
		request.Body.Read(buf)
		body := "Post route it is VT99 " + string(buf)
		return OK, body
	})
}

func NewEndpoints() *Endpoints {
	return &Endpoints{endpointsMap: map[string]map[HttpMethod]HandlerFunc{}}
}

// listen on a port using net or sth
// accept the connection on a port for an infinite time
// accept the connection
// create a goroutine for each connection
// handle the connection
// try removing the go keyword and see, ideally it should work in a sync manner

type temp struct {
	key string
	val int32
}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	defer listener.Close()
	log.Println("Server listening on :8000")

	endpoints = NewEndpoints()
	endpoints.registerEndpoint()

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

	values := endpoints.endpointsMap[request.URL.Path]
	// write the response to the connection
	log.Print(values, " values ")

	log.Print(endpoints.endpointsMap, " endpoints.endpointsMap ")

	if values == nil {
		fmt.Println("The value is Nil for HTTP Method: ", request.Method)
		WriteResponse(conn, MethodNotAllowed, "Method not allowed")
	}

	handler := values[HttpMethod(request.Method)]
	status, responseBody := handler(request)
	WriteResponse(conn, status, responseBody)
}

func WriteResponse(conn net.Conn, status StatusCode, response string) {

	response += "\r\n"
	headers := fmt.Sprintf("HTTP/1.1 %s\r\n", OK)
	headers += fmt.Sprintf("Content-Length: %d\r\n", len(response))
	headers += "Content-Type: text/plain\r\n\r\n"

	// write the headers and body to the connection
	fmt.Print(headers, "  headers ", response, " response ")
	conn.Write([]byte(headers))
	conn.Write([]byte(response))
}
