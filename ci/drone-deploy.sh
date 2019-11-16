#!/usr/bin/env bash

set -euxo pipefail

if [ -z "$DRONE" ]; then
    echo "this script is intended to be run only on drone" >&2;
    exit 1
fi

if [ ! -z "$DRONE_TAG" ]; then
	if [ "$(go run build.go version)" != "$DRONE_TAG" ]; then
        echo "octant version does not match tagged version!" >&2
        echo "octant version is $(go run build.go version)" >&2
        echo "tag is $DRONE_TAG" >&2
        exit 1
    fi

    goreleaser --rm-dist
fi
