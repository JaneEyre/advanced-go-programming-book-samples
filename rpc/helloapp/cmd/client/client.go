// cmd/client/main.go
package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	// Connect to the RPC server running on localhost:1234.
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Dialing error: %v", err)
	}
	// Ensure the connection is closed when the main function exits.
	defer client.Close()

	// The argument we are sending to the remote procedure.
	request := "Scalable Go RPC"
	// A variable to store the server's reply.
	var reply string

	// Perform the remote procedure call.
	// The first argument is the "Service.Method" name string.
	// The second argument is the request payload.
	// The third argument is a pointer to the variable where the reply will be stored.
	//err = client.Call("HelloService.Hello", request, &reply)
	err = client.Call("HiService.Hello", request, &reply)
	if err != nil {
		log.Fatalf("RPC call error: %v", err)
	}

	// If the call succeeds, print the response from the server.
	fmt.Println("Server Reply:", reply)
}
