#!/bin/bash

if [ -z "$TRAVIS" ]; then
    echo "this script is intended to be run only on travis" >&2;
    exit 1
fi

function goreleaser() {
    curl -sL https://git.io/goreleaser | bash
}


if [ ! -z "$TRAVIS_TAG" ]; then
	if [ "$(make version)" != "$TRAVIS_TAG" ]; then
        echo "hcli version does not match tagged version!" >&2
        echo "hcli version is $(make version)" >&2
        echo "tag is $TRAVIS_TAG" >&2
        exit 1
    fi

    goreleaser
fi
