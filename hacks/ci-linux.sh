#!/bin/sh

# install dependencies for CI

set -e

mkdir -p "$HOME"/bin

# protoc
curl -L https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip -o /tmp/protoc.zip
unzip /tmp/protoc.zip -d /tmp/protoc
cp /tmp/protoc/bin/protoc "$HOME"/bin/protoc
chmod 755 "$HOME"/bin/protoc
