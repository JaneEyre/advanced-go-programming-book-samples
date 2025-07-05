// Save this file as client.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	. "protobuf/hellopb"
)

func main() {
	// Dial the TCP server at the specified address.
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("net.Dial error: %v", err)
	}
	defer conn.Close()

	// Instead of rpc.NewClient, we use rpc.NewClientWithCodec to specify
	// that communication should be done using the JSON codec.
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// [NEW] proto
	// --- Changes start here ---

	// 1. Create a hellopb.String for the request
	requestArg := &String{
		Value: "Go proto client", // Set the actual string value inside the protobuf message
	}

	// 2. Declare a hellopb.String variable for the reply
	var reply String

	// 3. Call the remote method. Ensure the method name matches what's registered on the server.
	err = client.Call("HelloService.Hello", requestArg, &reply) // Pass the protobuf object
	if err != nil {
		log.Fatalf("RPC call error: %v", err)
	}

	// 4. Print the reply from the server. Access the Value field.
	fmt.Println("Server reply:", reply.GetValue())

	// --- Changes end here ---
}
