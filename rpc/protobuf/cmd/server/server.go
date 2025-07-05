// Save this file as server.go
package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	// The dot . before the package path tells the Go compiler
	// to bring all the exported identifiers (like String)
	// from the protobuf/hellopb package in server.go
	// without the need for the hellopb. prefix.
	. "protobuf/hellopb" // <--- Notice the dot here!
)

// HelloService is the struct that holds the methods we want to expose.
type HelloService struct{}

// Hello is the method we will expose via RPC.
// It must follow the rules for RPC methods in Go:
func (s *HelloService) Hello(request *String, reply *string) error {
	// [New] protobuf change
	*reply = "hello:" + request.GetValue()
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

		// This allows the server to handle multiple clients concurrently.
		// rpc.ServeCodec tells the RPC server to expect and respond with JSON.
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
