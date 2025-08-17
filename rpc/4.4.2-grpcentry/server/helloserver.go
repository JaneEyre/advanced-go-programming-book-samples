// server/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "hellogrpc/hellogrpc" // Import the generated code

	"google.golang.org/grpc"
)

// HelloServiceImpl implements the HelloServiceServer interface
type HelloServiceImpl struct {
	pb.UnimplementedHelloServiceServer
}

func (p *HelloServiceImpl) Hello(ctx context.Context, args *pb.String) (*pb.String, error) {
	log.Printf("Received: %v", args.GetValue())
	reply := &pb.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register our service implementation with the gRPC server
	pb.RegisterHelloServiceServer(grpcServer, new(HelloServiceImpl))

	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", lis.Addr())

	// Start serving gRPC requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
