//go:build tools
// +build tools

package discovery

import (
	_ "connectrpc.com/connect"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
