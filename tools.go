//go:build tools
// +build tools

package proto

import (
	_ "connectrpc.com/connect"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
