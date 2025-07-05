// reverse_rpc_backend.go (真实服务 - 通常在内网，主动连接代理)
package main

import (
	"log"
	"net"
	"net/rpc"
	"time"
)

// HelloService 是实际的 RPC 服务
type HelloService struct{}

func (s *HelloService) SayHello(name string, reply *string) error {
	*reply = "Hello, " + name + " from Backend (via reverse proxy)!"
	log.Printf("Backend (reverse) received SayHello for: %s", name)
	return nil
}

func main() {
	rpc.Register(new(HelloService)) // 注册 HelloService

	for {
		// **真实服务主动连接**到反向代理 (Proxy:1234)
		conn, err := net.Dial("tcp", "localhost:1234")
		if err != nil {
			log.Printf("Backend (reverse) dial proxy error: %v, retrying...", err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("Backend (reverse) connected to proxy: %s", conn.RemoteAddr())

		// **通过自己建立的连接提供 RPC 服务**
		// 代理会通过这条连接，将外部客户端的请求转发给这里的 HelloService
		rpc.ServeConn(conn)
		conn.Close() // 连接断开后重连
		log.Println("Backend (reverse) connection closed, retrying...")
	}
}
