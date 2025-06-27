// internal/helloservice/helloservice.go
package helloservice

import "net/rpc"

// HelloServiceName is the unique name for the RPC service.
// Defining it as a constant here ensures the client and server use the exact same string.
const HelloServiceName = "HelloService"

// HelloServiceInterface defines the set of methods that our service provides.
// Both the server-side implementation and the client-side stub can satisfy this interface.
type HelloServiceInterface interface {
	Hello(request string, reply *string) error
}

// RegisterHelloService is a helper function that wraps the standard rpc.RegisterName.
// It enforces that the registered service (svc) must satisfy the HelloServiceInterface,
// which prevents runtime errors by checking at compile time.
func RegisterHelloService(svc HelloServiceInterface) error {
	return rpc.RegisterName(HelloServiceName, svc)
}

// HelloService is the concrete implementation of the service on the server-side.
type HelloService struct{}

// The Hello method is the actual logic that runs on the server.
// It satisfies the HelloServiceInterface because its signature matches exactly.
func (p *HelloService) Hello(request string, reply *string) error {
	//*reply = HelloServiceName + "hello:" + request
	*reply = "ServiceName:" + HelloServiceName + " TO " + request
	return nil
}

// The Hello method is the actual logic that runs on the server.
// It satisfies the HelloServiceInterface because its signature matches exactly.
func (p *HelloService) Bye(request string, reply *string) error {
	//*reply = HelloServiceName + "hello:" + request
	*reply = "ServiceName:" + HelloServiceName + " 886 " + request
	return nil
}
