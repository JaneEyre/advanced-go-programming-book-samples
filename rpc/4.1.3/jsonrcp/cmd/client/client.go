// Save this file as client.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// Dial the TCP server at the specified address.
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("net.Dial error: %v", err)
	}
	// Defer closing the connection to ensure it's cleaned up.
	defer conn.Close()

	// Create a new RPC client.
	// Instead of rpc.NewClient, we use rpc.NewClientWithCodec to specify
	// that communication should be done using the JSON codec.
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// Declare a variable to store the server's reply.
	var reply string

	// Call the remote method "HelloService.Hello" with the argument "Go client".
	// The client handles marshalling the request into JSON, sending it,
	// and unmarshalling the JSON response into the 'reply' variable.
	err = client.Call("HelloService.Hello", "Go client", &reply)
	if err != nil {
		log.Fatalf("RPC call error: %v", err)
	}

	// Print the reply from the server.
	fmt.Println("Server reply:", reply)
}

