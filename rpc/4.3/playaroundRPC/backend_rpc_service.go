// backend_rpc_service.go
package main

import (
	"log"
	"net"
	"net/rpc"
)

// HelloService 是实际的 RPC 服务
type HelloService struct{}

func (s *HelloService) SayHello(name string, reply *string) error {
	*reply = "Hello, " + name + " from Backend!"
	log.Printf("Backend received SayHello for: %s", name)
	return nil
}

func main() {
	rpc.Register(new(HelloService))

	// 真实服务**监听**一个端口，等待连接
	listener, err := net.Listen("tcp", ":8000") // 后端服务监听在 8000
	if err != nil {
		log.Fatalf("Backend listen error: %v", err)
	}
	defer listener.Close()
	log.Println("Backend RPC Service listening on :8000")

	for {
		conn, err := listener.Accept() // 接受连接
		if err != nil {
			log.Printf("Backend accept error: %v", err)
			continue
		}
		go rpc.ServeConn(conn) // 通过连接提供 RPC 服务
	}
}
