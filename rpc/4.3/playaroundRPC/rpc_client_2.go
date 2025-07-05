// rpc_client.go
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func main() {
	// Client connects to the proxy's client-facing address
	client, err := rpc.Dial("tcp", "localhost:5678") // Connect to proxy's client port
	if err != nil {
		log.Fatalf("Client dial proxy error: %v", err)
	}
	defer client.Close()

	log.Println("Client connected to proxy, making RPC call...")

	var reply string
	// Call the ProxyService's SayHello method, which will be forwarded
	err = client.Call("ProxyService.SayHello", "Go Developer (External)", &reply)
	if err != nil {
		log.Fatalf("Client RPC call error: %v", err)
	}
	fmt.Printf("Client received reply: %s\n", reply)

	// Make another call to demonstrate routing to a potentially different backend if multiple exist
	time.Sleep(time.Millisecond * 500)
	err = client.Call("ProxyService.SayHello", "Another User", &reply)
	if err != nil {
		log.Fatalf("Client RPC call error (second): %v", err)
	}
	fmt.Printf("Client received second reply: %s\n", reply)

	time.Sleep(time.Second) // Keep client alive briefly for logs
}
