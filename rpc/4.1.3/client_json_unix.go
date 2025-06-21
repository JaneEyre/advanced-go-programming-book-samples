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

// GreeterServiceClient (reused, same as in Example 2)
type GreeterServiceClient struct {
	client *rpc.Client
}

// NewGreeterServiceClient creates a new RPC client connection for Unix Domain Socket
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
	err := g.client.Call("GreeterService.SayHello", name, &reply)
	return reply, err
}

func main() {
	socketPath := "/tmp/greeter.sock"
	fmt.Printf("Connecting to RPC Server via Unix Domain Socket at %s...\n", socketPath)

	client, err := NewGreeterServiceClient("unix", socketPath) // Dial using "unix" protocol
	if err != nil {
		log.Fatalf("Error connecting to RPC server: %v", err)
	}
	defer client.client.Close()

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

		reply, err := client.SayHello(name)
		if err != nil {
			log.Printf("Error calling remote method: %v", err)
			fmt.Println("Lost connection or RPC error. Exiting.")
			return
		}
		fmt.Printf("Client received reply: '%s'\n", reply)
	}
}
