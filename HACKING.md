# Hacking

## Requirements

* [Go 1.13 or above](https://golang.org/dl/)
* [node 10.15.0 or above](https://nodejs.org/en/)
* [npm 6.4.1 or above](https://www.npmjs.com/get-npm)
* [rice](https://github.com/GeertJohan/go.rice) - packaging web assets into a binary
* [mockgen](https://github.com/golang/mock) - generating go files used for testing
* [protoc](https://github.com/golang/protobuf) - generate go code compatible with gRPC

## Quick Start

    git clone git@github.com:vmware-tanzu/octant.git
    cd octant
    go run build.go go-install  # install Go dependencies.
    go run build.go ci-quick    # build UI, generate UI files, and create octant binary.
    ./build/octant   # run the Octant binary you just built

## Testing

We generally require tests be added for all but the most trivial of changes. You can run govet and the tests using the commands below:

    go run build.go vet
    go run build.go test

## Developing

When making changes to the frontend it can be helpful to have those changes trigger rebuilding the UI. Octant provides a short cut
using:

    go run build.go serve

The `serve` command starts two processes. The first is an alias for `npm run start` and will listen for changes and rebuild the UI.
The UI server will launch on `http://localhost:4200`.

The second, is an alias for `go run ./cmd/octant/main.go` but with useful environment variables already set, `OCTANT_PROXY_FRONTEND` which will reverse proxy to the Angular service and `OCTANT_DISABLE_OPEN_BROWSER` which prevents Octant from attempting to start the default system browser. The Octant server will launch on `http://localhost:7777`.

## Before Your Pull Request

When you are ready to create your pull request, we recommend running `go run build.go ci`.

This command will run our linting tools and test suite as well as produce a release binary that you can use to do a final
manual test of your changes.
