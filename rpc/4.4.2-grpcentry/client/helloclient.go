// client/main.go
package main

import (
	"context"
	"log"
	"time"

	pb "hellogrpc/hellogrpc" // Import the generated code

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:1234"
)

func main() {
	// Establish a connection to the gRPC server.
	// We use grpc.WithInsecure() for this example since we aren't using TLS.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a client stub for the HelloService
	client := pb.NewHelloServiceClient(conn)

	// Set a timeout for the RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the Hello method on the server
	r, err := client.Hello(ctx, &pb.String{Value: "world"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetValue())
}
