// reverse_rpc_proxy.go (代理端 - 外部可访问)
package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	// **反向代理监听**一个端口 (例如:1234)，等待**真实服务**来连接
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Reverse proxy listen error: %v", err)
	}
	defer listener.Close()
	log.Println("Reverse RPC Proxy listening on :1234")

	for {
		// 代理接受**真实服务**发来的连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Reverse proxy accept error: %v", err)
			continue
		}
		log.Printf("Reverse proxy accepted connection from %s (this is the Backend service)", conn.RemoteAddr())

		// **通过这条连接提供 RPC 服务** (给外部的客户端调用)
		// 注意：这里的 rpc.ServeConn(conn) 是让代理成为一个 RPC Server，
		// 外部客户端会把代理当成 RPC Server 来调用。
		// 但实际上，处理这些调用的逻辑是由**连上来的真实服务**通过这条连接反向传递完成的。
		// 在Go的net/rpc体系里，ServeConn就是将当前conn转换为一个RPC会话。
		// 这里的复杂性在于，ServeConn既可以接受客户端请求，也可以作为服务端处理请求。
		// 在反向RPC里，它接受外部客户端的请求，并通过conn转交给内部的真实服务。
		rpc.ServeConn(conn) // 代理通过这条连接，**为外部客户端提供 RPC 服务**
		conn.Close()        // 连接断开后关闭
	}
}
