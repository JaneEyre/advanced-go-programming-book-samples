cd /home/jiaguo/learnspace/advanced-go-programming-book-samples/<yourexamplefolder>

1) make sure protoc-gen-go-grpc is correctly installed
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3

export PATH="$PATH:$(go env GOPATH)/bin"

2) 
Add go_package to the .proto file , make sure to add content below 
'''
syntax = "proto3";
// Defines the Go package path for the generated code.
// must have a period in front, IF it means current folder
option go_package = ".hellogrpc";

3) code to generate proto service interface code
protoc --go_out=. --go-grpc_out=. ./hello.proto
go mod init hellogrpc // the mod to init is the name in generated pb.go code , NOT the one in orginal protobuf file
go mod tidy


4) in server/client implementation
4.1) follow structure
my-grpc-project/
├── go.mod
├── hello.proto
├── client/
│   └── main.go           # Contains package main
└── server/
    └── main.go           # Contains package main
└── hellogrpc/
    ├── hello.pb.go       # Contains package hellogrpc
    └── hello_grpc.pb.go  # Contains package hellogrpc

in 2 main.go , your package import path should look like:
'''
package main

import (
	pb "hellogrpc/hellogrpc" // Import the generated code
)
'''