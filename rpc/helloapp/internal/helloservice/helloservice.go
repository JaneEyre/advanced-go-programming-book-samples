// internal/helloservice/service.go
package helloservice

// This file defines the shared RPC service.
// By placing it in an 'internal' package, we ensure it can only be imported
// by other packages within our 'my-rpc-app' module.

// HelloService is the type that provides the RPC methods.
// It's an empty struct because its purpose is to be a container for the methods.
type HelloService struct{}

// Hello is the method that will be exposed via RPC.
// It follows the rules for Go's net/rpc package:
// - It is an exported method (starts with a capital letter).
// - It has two arguments, both of which are exported or built-in types.
// - The second argument is a pointer.
// - It returns a single value of type 'error'.
func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}
