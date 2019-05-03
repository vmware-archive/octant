#!/bin/bash

# install dependencies for CI

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source ${DIR}/functions.sh

mkdir -p "$HOME"/bin

install_protoc_linux
