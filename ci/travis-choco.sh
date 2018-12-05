#!/bin/bash

set -euxo pipefail

if [ -z "$TRAVIS" ]; then
    echo "this script is intended to be run only on travis" >&2;
    exit 1
fi

if [[ $TRAVIS_OS_NAME == 'windows' ]]; then
    choco pack $TRAVIS_BUILD_DIR/choco/hcli.nuspec
    choco install $TRAVIS_BUILD_DIR/hcli.$(make version | cut -c2- ).nupkg
    hcli version
fi

