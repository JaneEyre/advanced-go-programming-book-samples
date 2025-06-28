// cmd/server/main.go
package main

import (
	"hellointerface/internal/helloservice"
	"log"
	"net"
	"net/rpc"
)

func main() {
	helloservice.RegisterHelloService(new(helloservice.HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	log.Println("Server listening on port 1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
			continue
		}
		// *** ADD THIS LINE ***
		// Log the address of the client that just connected.
		log.Printf("Accepted connection from: %s", conn.RemoteAddr().String())

		go rpc.ServeConn(conn)
	}
}
