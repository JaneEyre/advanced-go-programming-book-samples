// client/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "stream/stream" // Update with your module path

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:1234"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	// Call the Channel method to get the bidirectional stream object
	stream, err := client.Channel(context.Background())
	if err != nil {
		log.Fatalf("could not get stream: %v", err)
	}

	// Goroutine for sending data to the server
	go func() {
		for i := 0; ; i++ {
			msg := fmt.Sprintf("hi from client: %d", i)
			if err := stream.Send(&pb.String{Value: msg}); err != nil {
				log.Fatalf("failed to send: %v", err)
			}
			log.Printf("Client sent: %s", msg)
			time.Sleep(time.Second)
		}
	}()

	// Main loop for receiving data from the server
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break // Server closed the stream
			}
			log.Fatalf("failed to receive: %v", err)
		}
		fmt.Printf("Received from server: %s\n", reply.GetValue())
	}
}
