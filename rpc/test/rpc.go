package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    // Just to use the import and avoid unused import error
    _ = fmt.Sprintf("Checking rpc.Codec: %T", (rpc.Codec)(nil))
    fmt.Println("net/rpc imported successfully.")
}
