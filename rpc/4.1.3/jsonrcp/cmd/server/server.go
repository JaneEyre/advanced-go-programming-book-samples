// Save this file as server.go
package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// HelloService is the struct that holds the methods we want to expose.
type HelloService struct{}

// Hello is the method we will expose via RPC.
// It must follow the rules for RPC methods in Go:
// - It must be an exported method of an exported type.
// - It must accept two arguments, both exported (or built-in) types.
// - The second argument must be a pointer.
// - It must have a return type of error.
func (s *HelloService) Hello(request string, reply *string) error {
	// The core logic of the method.
	*reply = "hello:" + request
	return nil
}

func main() {
	// Register the HelloService with the RPC system.
	// This makes its public methods available to be called.
	err := rpc.RegisterName("HelloService", new(HelloService))
	if err != nil {
		log.Fatalf("Failed to register RPC service: %v", err)
	}

	// Listen for incoming TCP connections on port 1234.
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("ListenTCP error: %v", err)
	}
	// Defer closing the listener to ensure it's cleaned up on exit.
	defer listener.Close()
	log.Println("Server is listening on port 1234")

	// Loop indefinitely to accept and handle new connections.
	for {
		// Accept blocks until a new connection is made.
		conn, err := listener.Accept()
		if err != nil {
			// If an error occurs, log it and continue to the next iteration.
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle each client connection in a new goroutine.
		// This allows the server to handle multiple clients concurrently.
		// We use rpc.ServeCodec with the JSON codec to handle the communication.
		// This tells the RPC server to expect and respond with JSON.
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

