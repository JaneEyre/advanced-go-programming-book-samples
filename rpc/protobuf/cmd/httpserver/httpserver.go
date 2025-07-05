// Save this file as server.go
package main

import (
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// HelloService is the struct that holds the methods we want to expose.
type HelloService struct{}

// Hello is the method we will expose via RPC.
// It must follow the rules for RPC methods in Go:
// - It must be an exported method of an exported type.
// - It must accept two arguments, both exported (or built-in) types.
// - The second argument must be a pointer.
// - It must have a return type of error.
func (s *HelloService) Hello(request string, reply *string) error {
	// The core logic of the method.
	*reply = "hello:" + request
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {

		// --- MORE DETAILED LOGGING ---
		// Get the User-Agent header from the request
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			userAgent = "Unknown"
		}
		// Log the client's IP, Method, and User-Agent
		log.Printf("Accepted connection from: %s, Method: %s, User-Agent: %s", r.RemoteAddr, r.Method, userAgent)
		// --- END OF LOGGING ---

		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})
	log.Println("HTTP JSON-RPC server listening on :1234/jsonrpc")

	http.ListenAndServe(":1234", nil)
}
