PROTOTOOL_VERSION="1.8.0"
GRPC_GATEWAY_VERSION="1.11.1"
PROTOC_GEN_GO_VERSION="1.3.2"

# get protoc-gen-grpc-gateway and protoc-gen-go and add them to $PATH
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION}
go get github.com/golang/protobuf/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}
