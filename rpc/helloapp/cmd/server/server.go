// cmd/server/main.go
package main

import (
	"helloapp/internal/helloservice"
	"log"
	"net"
	"net/rpc"
	// Import the shared service package.
	// Replace 'my-rpc-app' with your actual go module name from go.mod.
)

func main() {
	// Create a new instance of our HelloService and register it with the RPC system.
	// This makes its methods available for remote calls under the name "HelloService".
	//err := rpc.RegisterName("HelloService", new(helloservice.HelloService))
	err := rpc.RegisterName("HiService", new(helloservice.HelloService))
	if err != nil {
		log.Fatalf("Failed to register RPC service: %v", err)
	}

	// Listen for incoming TCP connections on port 1234.
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("ListenTCP error: %v", err)
	}
	// Ensure the listener is closed when the main function exits.
	defer listener.Close()

	log.Println("RPC Server is listening on port 1234")

	// Enter an infinite loop to continuously accept new client connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			// Continue to the next iteration to accept other connections.
			continue
		}

		// Handle each client connection in its own goroutine.
		// This allows the server to serve multiple clients concurrently without blocking.
		go rpc.ServeConn(conn)
	}
}
