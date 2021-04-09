#!/bin/sh
# generate golang for protobuf

# get directory of this script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
OCTANT_ROOT=${DIR}/../../..
MODULE="github.com/vmware-tanzu/octant/pkg/plugin/dashboard"

protoc -I${OCTANT_ROOT}/vendor -I${OCTANT_ROOT} -I${DIR} --go_out=plugins=grpc:${DIR} --go_opt=module=${MODULE} ${DIR}/dashboard.proto
