package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
	"time"
)

// GreeterServiceClient is a client-side wrapper for the RPC service
type GreeterServiceClient struct {
	client *rpc.Client
}

// NewGreeterServiceClient creates a new RPC client connection with TCP and custom codec
func NewGreeterServiceClient(network, address string) (*GreeterServiceClient, error) {
	conn, err := net.DialTimeout(network, address, time.Second*5) // Add connection timeout
	if err != nil {
		return nil, err
	}
	// Use custom JsonCodec for the client
	client := rpc.NewClientWithCodec(NewJsonCodec(conn))
	return &GreeterServiceClient{client: client}, nil
}

// SayHello makes an RPC call to the remote SayHello method
func (g *GreeterServiceClient) SayHello(name string) (string, error) {
	var reply string
	// Call the remote method: "Service.Method", args, &reply
	// "GreeterService" is the name registered on the server (from type name)
	err := g.client.Call("GreeterService.SayHello", name, &reply)
	return reply, err
}

func main() {
	fmt.Println("Connecting to RPC Server at localhost:1234...")
	client, err := NewGreeterServiceClient("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Error connecting to RPC server: %v", err)
	}
	defer client.client.Close() // Close the RPC client connection when main exits

	fmt.Println("Connected. Type your name and press Enter to get a greeting. Type 'exit' to quit.")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your name: ")
		input, _ := reader.ReadString('\n')
		name := strings.TrimSpace(input)

		if strings.ToLower(name) == "exit" {
			fmt.Println("Exiting.")
			return
		}

		// Make the RPC call
		reply, err := client.SayHello(name)
		if err != nil {
			log.Printf("Error calling remote method: %v", err)
			// Decide if you want to exit or continue trying
			fmt.Println("Lost connection or RPC error. Exiting.")
			return
		}

		fmt.Printf("Received RPC reply: '%s'\n", reply)
	}
}
