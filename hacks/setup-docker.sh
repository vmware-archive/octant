#!/bin/bash

set -ex -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

source ${DIR}/functions.sh

function init() {
    apt update
    apt install -y zip apt-utils
}

init
install_mockgen
install_rice
install_npm
install_protoc_linux

make generate web-build octant-docker
