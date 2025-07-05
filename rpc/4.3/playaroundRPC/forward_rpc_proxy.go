// forward_rpc_proxy.go
package main

import (
	"io"
	"log"
	"net"
)

func handleForwardProxyConn(clientConn net.Conn) {
	defer clientConn.Close()

	// **代理主动连接**到后端 RPC 服务 (Backend:8000)
	backendConn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Printf("Proxy could not connect to backend: %v", err)
		return
	}
	defer backendConn.Close()

	// 双向转发数据
	go io.Copy(backendConn, clientConn) // 客户端请求 -> 后端服务
	io.Copy(clientConn, backendConn)    // 后端服务响应 -> 客户端
}

func main() {
	// **代理监听**一个端口 (例如:1234)，等待客户端连接
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Proxy listen error: %v", err)
	}
	defer listener.Close()
	log.Println("Forward RPC Proxy listening on :1234")

	for {
		conn, err := listener.Accept() // 接受客户端连接
		if err != nil {
			log.Printf("Proxy accept error: %v", err)
			continue
		}
		go handleForwardProxyConn(conn) // 为每个客户端连接启动一个处理协程
	}
}
