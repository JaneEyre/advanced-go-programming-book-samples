// server/main.go
package main

import (
	"context"
	"io"
	"log"
	"net"

	pb "stream/stream" // Update with your module path

	"google.golang.org/grpc"
)

// HelloServiceImpl implements the HelloServiceServer interface
type HelloServiceImpl struct {
	pb.UnimplementedHelloServiceServer
}

func (p *HelloServiceImpl) Hello(ctx context.Context, args *pb.String) (*pb.String, error) {
	log.Printf("Received Hello request: %v", args.GetValue())
	return &pb.String{Value: "hello:" + args.GetValue()}, nil
}

// Channel implements the bidirectional streaming RPC.
func (p *HelloServiceImpl) Channel(stream pb.HelloService_ChannelServer) error {
	for {
		// Receive data from the client's stream
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil // Client stream closed
			}
			return err
		}

		// Process the received data and create a reply
		reply := &pb.String{Value: "hello:" + args.GetValue()}

		// Send the reply back to the client
		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &HelloServiceImpl{})

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
