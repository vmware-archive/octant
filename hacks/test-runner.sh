#!/bin/bash

set -euxo pipefail

FRONTEND="web" # Specify multiple directory names separated by space
BUILD_FRONTEND=0
IGNORE_FILES=$(ls -p | grep -v /)

is_updated() {
    GIT_DIFF=$(git --no-pager diff --name-only HEAD~1 | sort -u | awk 'BEGIN {FS="/"} {print $1}' | uniq); 
    export CHANGED=$GIT_DIFF
}

should_build() {
    WATCH=($(echo "$FRONTEND" | tr ' ' '\n'))
    is_updated
    MODIFIED=($(echo "$CHANGED"))
    for j in "${WATCH[@]}"
    do
        for k in "${MODIFIED[@]}"
	do
            if [[ $j = $k ]]; then
                BUILD_FRONTEND=1
	    fi
	done
    done
}

if [[ $TRAVIS_BRANCH == 'master' ]]; then
    make ci
else
  should_build
  if [[ $BUILD_FRONTEND == 0 ]]; then
      # Backend tests only
      make test vet octant-dev
  else
    make web-test web-build
  fi
fi

exit $?
