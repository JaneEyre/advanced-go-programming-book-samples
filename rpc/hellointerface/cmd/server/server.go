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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
