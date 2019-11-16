#!/usr/bin/env bash

# This script makes sure that the checked-in mock files ("fake" packages) are
# up-to-date.

go run build.go go-install
find pkg -type d -name fake -prune -exec rm -rf {} \;
find internal -type d -name fake -prune -exec rm -rf {} \;
go run build.go generate
diff=$(git status --porcelain | grep fake)

if [ ! -z "$diff" ]; then
    echo "The generated mock files are not up-to-date" >&2
    echo "You can regenerate them with:" >&2
    echo "go run build.go go-install && go run build.go generate" >&2
    exit 1
fi
