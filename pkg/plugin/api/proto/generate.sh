#!/bin/sh
# generate golang for protobuf

protoc -I$GOPATH/src/github.com/vmware/octant/vendor -I$GOPATH/src/github.com/vmware/octant -I. --go_out=plugins=grpc:. dashboard.proto
