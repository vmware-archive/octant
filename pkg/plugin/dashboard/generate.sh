#!/bin/sh
# generate golang for protobuf

MODULE="github.com/vmware-tanzu/octant/pkg/plugin/dashboard"

protoc --go_out=. --go-grpc_out=paths=source_relative:. --go_opt=module=${MODULE} dashboard.proto
