package main

import (
	"context"
	"log"
	"os"
	"time"

	// Import the generated Go code
	pb "gRPC/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	// We use WithTransportCredentials and insecure.NewCredentials() because
	// we are not using SSL/TLS in this example.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client stub from the connection.
	c := pb.NewGreeterClient(conn)

	// Determine the name to greet from command-line arguments.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// Contact the server and print out its response.
	// We create a context with a 1-second deadline for the RPC call.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make the RPC call to the SayHello method.
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// Print the 'message' field from the HelloReply struct.
	log.Printf("Greeting: %s", r.GetMessage())
}
