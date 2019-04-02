package proto

//go:generate protoc -I$GOPATH/src/github.com/heptio/developer-dash/vendor -I$GOPATH/src/github.com/heptio/developer-dash -I. --go_out=plugins=grpc:. dashboard.proto
