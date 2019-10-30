# Hacking

## Requirements

* [Go 1.13 or above](https://golang.org/dl/)
* [node 10.15.0 or above](https://nodejs.org/en/)
* [npm 6.4.1 or above](https://www.npmjs.com/get-npm)
* [rice](https://github.com/GeertJohan/go.rice) - packaging web assets into a binary
* [mockgen](https://github.com/golang/mock) - generating go files used for testing
* [protoc](https://github.com/golang/protobuf) - generate go code compatible with gRPC

## Quick Start

```sh
git clone git@github.com:vmware-tanzu/octant.git
cd octant

# Manually install required Go tools as listed above, by following these instructions:
# - https://github.com/GeertJohan/go.rice#installation and
#   `export PATH="$PATH:${GOPATH}/bin"`
# - https://github.com/golang/mock#installation
# - https://github.com/golang/protobuf#installation

make go-install  # install Go dependencies.
make web-deps    # install npm dependencies (one-time step, calls `npm ci`;
                 # alternatively use `(cd web && npm install)` to avoid
                 # redownloading all modules)
make ci-quick    # build UI, generate UI files, and create octant binary.
./build/octant   # run the Octant binary you just built
```

## Testing

We generally require tests be added for all but the most trivial of changes. You can run govet and the tests using the commands below:

```sh
make vet
make test
```

## Frontend

When making changes to the frontend it can be helpful to have those changes trigger rebuilding the UI.
The Octant makefile provides `make ui-client` which is an alias for `npm run start` and will listen for changes and rebuild the UI.
By default this will launch on `http://localhost:4200`.

## Backend

When you are making changes to the backend you can take advantage of Go's fast compile time to build and run
Octant in a single step. The Octant makefile provides `make ui-server` which is an alias for `go run`. Unlike the
alias for the frontend, this does not listen for changes and does require you to stop the command and re-run it after
saving your changes.

If working on the frontend, you may want to set up a reverse proxy to the Angular services running on `localhost:4200`.
To set this up, set `OCTANT_PROXY_FRONTEND` environment variable with the location of the frontend.
(e.g. `http://localhost:4200`).

## Before Your Pull Request

When you are ready to create your pull request, we recommend running `make ci`.

This command will run our linting tools and test suite as well as produce a release binary that you can use to do a final
manual test of your changes.
