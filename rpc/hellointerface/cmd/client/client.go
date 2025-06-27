// cmd/client/main.go
package main

import (
	"fmt"
	"hellointerface/internal/helloservice"
	"log"
	"net/rpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

var _ helloservice.HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(helloservice.HelloServiceName+".Hello", request, reply)
}

func (p *HelloServiceClient) Bye(request string, reply *string) error {
	return p.Client.Call(helloservice.HelloServiceName+".Bye", request, reply)
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Hello("Client", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)

	err = client.Bye("Client2 ", &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)
}
