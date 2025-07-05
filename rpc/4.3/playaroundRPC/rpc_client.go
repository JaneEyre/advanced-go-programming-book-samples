// rpc_client.go
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

func main() {
	// 客户端连接到**代理**的地址 (无论是正向还是反向，都是 :1234)
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatalf("Client dial proxy error: %v", err)
	}
	defer client.Close()

	log.Println("Client connected to proxy, making RPC call...")

	var reply string
	err = client.Call("HelloService.SayHello", "Go Developer", &reply) // 调用 RPC 方法
	if err != nil {
		log.Fatalf("Client RPC call error: %v", err)
	}
	fmt.Printf("Client received reply: %s\n", reply)

	time.Sleep(time.Second)
}
