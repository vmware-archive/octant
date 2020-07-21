// +build tools

package tools

import (
	_ "github.com/GeertJohan/go.rice/rice"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "golang.org/x/tools/cmd/goimports"

	_ "k8s.io/client-go/dynamic/dynamicinformer" // Used for generated fakes
)
