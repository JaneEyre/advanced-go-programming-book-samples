package main

import (
	"context"
	"log"
	"net"

	// Import the generated Go code from our .proto file
	pb "gRPC/hello"

	"google.golang.org/grpc"
)

// server is used to implement the GreeterServer interface from our .proto file.
// It must have a method with the exact signature of SayHello.
type server struct {
	// This is a requirement for forward compatibility in gRPC.
	// It ensures that if the .proto file is updated with new methods,
	// this server will still compile.
	pb.UnimplementedGreeterServer
}

// SayHello implements the SayHello RPC method.
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// Log the name received in the request.
	log.Printf("Received name: %v", in.GetName())
	// Construct the reply and send it back.
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	// Create a TCP listener on port 50051. gRPC typically uses ports in this range.
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server instance.
	s := grpc.NewServer()

	// Register our service implementation with the gRPC server.
	// This connects our 'server' struct logic to the 'Greeter' service.
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	// Start the server and block until it's terminated.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
