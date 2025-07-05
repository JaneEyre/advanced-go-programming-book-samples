// reverse_rpc_proxy.go
package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// --- Internal Proxy Structure ---

// BackendInfo holds an RPC client to a connected backend
type BackendInfo struct {
	ID        string      // Unique ID for this backend connection
	RPCClient *rpc.Client // RPC client to talk to this specific backend
	Conn      net.Conn    // The underlying net.Conn from the backend
	LastUsed  time.Time   // For simple load balancing (least recently used)
}

var (
	// Pool of active backend connections. In a real system, you'd manage this more robustly
	// with health checks, registration/deregistration, etc.
	backendPool  map[string]*BackendInfo
	poolMu       sync.RWMutex
	backendCount int // To generate unique IDs
)

func init() {
	backendPool = make(map[string]*BackendInfo)
}

// getBackend selects an available backend from the pool (simple least recently used)
func getBackend() (*BackendInfo, error) {
	poolMu.RLock()
	defer poolMu.RUnlock()

	if len(backendPool) == 0 {
		return nil, fmt.Errorf("no active backends available")
	}

	var leastRecentlyUsed *BackendInfo
	for _, backend := range backendPool {
		if leastRecentlyUsed == nil || backend.LastUsed.Before(leastRecentlyUsed.LastUsed) {
			leastRecentlyUsed = backend
		}
	}

	if leastRecentlyUsed != nil {
		leastRecentlyUsed.LastUsed = time.Now() // Update last used time
	}
	return leastRecentlyUsed, nil
}

// --- Proxy's RPC Service (what external clients call) ---

// ProxyService is the service exposed by the proxy to external clients
type ProxyService struct{}

// SayHello acts as a proxy for the actual SayHello method on the backend
func (p *ProxyService) SayHello(name string, reply *string) error {
	backend, err := getBackend()
	if err != nil {
		return err // No backend available
	}

	log.Printf("Proxy: Routing SayHello(%s) to backend %s", name, backend.ID)

	// Call the actual backend service using the RPC client stored for that backend
	err = backend.RPCClient.Call("HelloService.SayHello", name, reply)
	if err != nil {
		log.Printf("Proxy: Error calling backend %s: %v", backend.ID, err)
		// Consider removing this backend from the pool if it fails
	}
	return err
}

// --- Main Proxy Logic ---

func main() {
	// Register the ProxyService so external clients can call it
	rpc.Register(new(ProxyService))

	// Goroutine 1: Listener for BACKEND connections (from backend_rpc_service.go)
	go func() {
		backendListener, err := net.Listen("tcp", ":1234") // Proxy listens for backends on 1234
		if err != nil {
			log.Fatalf("Proxy backend listener error: %v", err)
		}
		defer backendListener.Close()
		log.Println("Proxy: Listening for backends on :1234")

		for {
			conn, err := backendListener.Accept() // Accepts connection from a backend
			if err != nil {
				log.Printf("Proxy: Backend accept error: %v", err)
				continue
			}

			// Handle each backend connection in a new goroutine
			go func(backendConn net.Conn) {
				id := fmt.Sprintf("backend-%d", backendCount)
				backendCount++
				log.Printf("Proxy: Accepted backend connection from %s, ID: %s", backendConn.RemoteAddr(), id)

				// Create an rpc.Client using this connection to talk to the backend
				backendRPCClient := rpc.NewClient(backendConn)
				// Note: rpc.NewClient blocks if there's no server on the other end,
				// but here we know the backend is calling ServeConn on its side.

				backendInfo := &BackendInfo{
					ID:        id,
					RPCClient: backendRPCClient,
					Conn:      backendConn,
					LastUsed:  time.Now(),
				}

				poolMu.Lock()
				backendPool[id] = backendInfo // Add to pool
				poolMu.Unlock()

				// Crucially, we need to monitor this connection. If the backend disconnects,
				// the RPCClient will fail. We must remove it from the pool.
				// A simple way is to try a dummy call or just let subsequent actual calls fail.
				// For robustness, you'd add health checks or listen for disconnects.
				// For now, if the connection dies, the next RPC call using this client will error.

				// Keep this goroutine alive as long as the backend connection is.
				// One way to manage this is to wait for the connection to be closed from the other end.
				// This loop here is mostly to keep the goroutine active and clean up on disconnect.
				buf := make([]byte, 1) // Small buffer for reading to detect disconnects
				for {
					backendConn.SetReadDeadline(time.Now().Add(time.Second * 5)) // Set a deadline for Read
					_, err := backendConn.Read(buf)                              // Try to read to detect if connection is still alive
					if err != nil {
						if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
							// Timeout, connection still seems fine, continue loop
							continue
						}
						// Real error or EOF (disconnect)
						log.Printf("Proxy: Backend %s connection error: %v. Removing from pool.", id, err)
						break // Exit loop, trigger deferred Close
					}
					// If we successfully read, the backend is sending us data unexpectedly,
					// or it's a KeepAlive, etc. For RPC, client doesn't send much raw data.
					// This part is less common in simple RPC, but necessary for robust monitoring.
				}

				// Clean up when connection closes
				poolMu.Lock()
				delete(backendPool, id)
				poolMu.Unlock()
				backendRPCClient.Close() // Close the RPC client
				log.Printf("Proxy: Backend %s disconnected. Removed from pool.", id)

			}(conn) // Pass the accepted connection to the goroutine
		}
	}()

	// Goroutine 2 (main goroutine): Listener for EXTERNAL CLIENT connections
	clientListener, err := net.Listen("tcp", ":5678") // Proxy listens for external clients on 5678
	if err != nil {
		log.Fatalf("Proxy client listener error: %v", err)
	}
	defer clientListener.Close()
	log.Println("Proxy: Listening for external clients on :5678")

	for {
		conn, err := clientListener.Accept() // Accepts connection from an external client
		if err != nil {
			log.Printf("Proxy: External client accept error: %v", err)
			continue
		}
		log.Printf("Proxy: Accepted external client connection from %s", conn.RemoteAddr())

		// Serve RPC requests from external clients on this connection
		// ProxyService methods will internally route to backends.
		go rpc.ServeConn(conn) // This blocks for the lifetime of the client's RPC session
	}
}
