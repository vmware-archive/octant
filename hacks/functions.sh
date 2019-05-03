#!/bin/bash

function install_mockgen() {
    go get golang.org/x/tools/go/packages
    go get github.com/golang/mock/gomock
    go install github.com/golang/mock/mockgen
}

function install_rice() {
    go get github.com/GeertJohan/go.rice/rice
}

function install_npm() {
   curl -sL https://deb.nodesource.com/setup_10.x | bash
   apt install nodejs
}

function install_protoc_linux() {
    curl -L https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip -o /tmp/protoc.zip
    unzip /tmp/protoc.zip -d /tmp/protoc
    mkdir -p "$HOME/bin"
    cp /tmp/protoc/bin/protoc "$HOME"/bin/protoc
    chmod 755 "$HOME"/bin/protoc
}

function install_protoc_macos() {
    curl -L https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-osx-x86_64.zip -o /tmp/protoc.zip
    unzip /tmp/protoc.zip -d /tmp/protoc
    mkdir -p "$HOME/bin"
    cp /tmp/protoc/bin/protoc "$HOME"/bin/protoc
    chmod 755 "$HOME"/bin/protoc
}
