package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// GreeterService (reused, same as in Example 2)
type GreeterService struct{}

func (g *GreeterService) SayHello(name string, reply *string) error {
	log.Printf("Server received: SayHello(%s)", name)
	*reply = fmt.Sprintf("Hello, %s! (via Unix Socket RPC)", name)
	return nil
}

func main() {
	socketPath := "/tmp/greeter.sock"

	// Ensure socket file does not exist to prevent "address already in use" error
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Error removing old socket: %v", err)
	}

	rpc.Register(new(GreeterService))

	// Listen using "unix" protocol
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	defer listener.Close()
	fmt.Printf("RPC server with Unix Domain Socket listening on %s\n", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		// Still use NewJsonCodec because it works with any io.ReadWriteCloser
		go rpc.ServeCodec(NewJsonCodec(conn))
	}
}
