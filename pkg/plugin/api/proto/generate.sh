#!/bin/sh
# generate golang for protobuf

# get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
OCTANT_ROOT=${DIR}/../../../..

protoc -I${OCTANT_ROOT}/vendor -I${OCTANT_ROOT} -I. --go_out=plugins=grpc:. dashboard_api.proto

