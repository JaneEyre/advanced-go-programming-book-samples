package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// GreeterService is the type that will provide RPC methods
type GreeterService struct{}

// SayHello is the RPC method that will be exposed
// It must follow the signature: (argsType, *replyType) error
func (g *GreeterService) SayHello(name string, reply *string) error {
	log.Printf("Server received: SayHello(%s)", name)
	*reply = fmt.Sprintf("Hello, %s! (via JSON RPC)", name)
	return nil // No error
}

func main() {
	rpc.Register(new(GreeterService)) // Registers the service using its type name

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	defer listener.Close()

	fmt.Println("RPC server with JSON Codec listening on :1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		// Use the custom JsonCodec for each connection
		go rpc.ServeCodec(NewJsonCodec(conn))
	}
}
