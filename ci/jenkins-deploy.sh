#!/bin/bash

set -euxo pipefail

if [ -z "$JENKINS_HOME" ]; then
    echo "this script is intended to be run only on jenkins" >&2;
    exit 1
fi

function goreleaser() {
    curl -sL https://git.io/goreleaser | bash
}


# GIT_TAG_NAME comes from Git Tag Message Plugin
# Confirms semver in Makefile is consistent with pushed tag
if [ ! -z "$GIT_TAG_NAME" ]; then
	if [ "$(make version)" != "$GIT_TAG_NAME" ]; then
        echo "version does not match tagged version!" >&2
        echo "version is $(make version)" >&2
        echo "tag is $GIT_TAG_NAME" >&2
        exit 1
    fi

    goreleaser release --skip-publish --rm-dist
fi
