#!/usr/bin/env bash

set -euxo pipefail

if [ -z "$GITHUB_TAG" ]; then
    echo "this script is intended to be run only on github actions" >&2;
    exit 1
fi

if [ ! -z "$GITHUB_TAG" ]; then
	if [ "$(go run build.go version)" != "$GITHUB_TAG" ]; then
        echo "octant version does not match tagged version!" >&2
        echo "octant version is $(go run build.go version)" >&2
        echo "tag is $GITHUB_TAG" >&2
        exit 1
    fi

    exit 0
fi
