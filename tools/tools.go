// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"

	_ "k8s.io/client-go/dynamic/dynamicinformer" // Used for generated fakes
)
