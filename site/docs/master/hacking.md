# Hacking

## Requirements

* [Go 1.13 or above](https://golang.org/dl/)
* [node 10.15.0 or above](https://nodejs.org/en/)
* [npm 6.4.1 or above](https://www.npmjs.com/get-npm)
* [rice](https://github.com/GeertJohan/go.rice) - packaging web assets into a binary
* [mockgen](https://github.com/golang/mock) - generating go files used for testing
* [protoc](https://github.com/golang/protobuf) - generate go code compatible with gRPC

These build tools can be installed with `go run build.go go-install`.

A development binary can be built by `go run build.go build`.

For UI changes, see the [README]({{ site.gh_repo }}/tree/master/web) located in `web/`.

If Docker and [Drone](/docs/drone) are installed, tests and build steps can run in a containerized environment.

### Developer Variables

These variables are for aiding in the development of Octant. They should not be needed for normal operation of Octant. These variables are subject to change and/or removal at anytime and will not follow a deprecation process.

* `OCTANT_DISABLE_OPEN_BROWSER` - set to a true value if you do not want the browser launched when the dashboard starts up.
* `OCTANT_DISABLE_CLUSTER_OVERVIEW` - set to a true value if you do not want the cluster overview module to load.
* `OCTANT_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)
* `OCTANT_ACCEPTED_HOSTS` - set to comma-separated string of hosts to be accepted. (e.g. `demo.octant.example.com,awesome.octant.zr`)
* `OCTANT_LOCAL_CONTENT` - set to a directory and dash will serve content responses from here. An example directory lives in `examples/content`

## e2e Testing

Cypress will load the dashboard from port 7777. Navigate to `web/` then install the Cypress binary with:

```sh
npm install cypress --save-dev
```

Run the test from the command line with the option of specifying a browser or electron:

```sh
$(npm bin)/cypress run -b chrome
```

Starts the interactive launcher to load tests in `/cypress`.

```sh
$(npm bin)/cypress open
```

## Quick Start

```sh
git clone git@github.com:vmware-tanzu/octant.git
cd octant

# Manually install required Go tools as listed above, by following these instructions:
# - https://github.com/GeertJohan/go.rice#installation and
#   `export PATH="$PATH:${GOPATH}/bin"`
# - https://github.com/golang/mock#installation
# - https://github.com/golang/protobuf#installation

go run build.go go-install  # install Go dependencies.
go run build.go web-deps    # install npm dependencies (one-time step, calls `npm ci`;
                            # alternatively use `(cd web && npm install)` to avoid
                            # redownloading all modules)
go run build.go ci-quick    # build UI, generate UI files, and create octant binary.
./build/octant              # run the Octant binary you just built
```

## Testing

We generally require tests be added for all but the most trivial of changes. You can run govet and the tests using the commands below:

```sh
go run build.go vet
go run build.go test
```

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
