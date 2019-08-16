# Getting Started

## Environment Variables

Octant is configurable through environment variables defined at runtime.

* `KUBECONFIG` - set to non-empty location if you want to set KUBECONFIG with an environment variable.
* `OCTANT_DISABLE_OPEN_BROWSER` - set to a non-empty value if you don't the browser launched when the dashboard start up.
* `OCTANT_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)
* `OCTANT_VERBOSE_CACHE` - set to a non-empty value to view cache actions
* `OCTANT_LOCAL_CONTENT` - set to a directory and dash will serve content responses from here. An example directory lives in `examples/content`
* `OCTANT_PLUGIN_PATH` - add a plugin directory or multiple directories separated by `:`. Plugins will load by default from `$HOME/.config/octant/plugins`

**Note:** If using [fish shell](https://fishshell.com), tilde expansion may not occur when using `env` to set environment variables.

## Setting Up a Development Environment

* [Go 1.12 or above](https://golang.org/dl/)
* [node 10.15.0 or above](https://nodejs.org/en/)
* [npm 6.4.1 or above](https://www.npmjs.com/get-npm)
* [rice](https://github.com/GeertJohan/go.rice) - packaging web assets into a binary
* [mockgen](https://github.com/golang/mock) - generating go files used for testing
* [protoc](https://github.com/golang/protobuf) - generate go code compatible with gRPC

These build tools can be installed via Makefile with `make go-install`.

A development binary can be built by `make octant-dev`.

For UI changes, see the [README](/web/README.md) located in `web/`.

If Docker and [Drone](/docs/drone.md) are installed, tests and build steps can run in a containerized environment.

## e2e Testing

Cypress will load the dashboard from port 7777. Navigate to `web/` then install the Cypress binary `npm install cypress --save-dev`.

Run the test from the command line with the option of specifying a browser or electron:

`$(npm bin)/cypress run -b chrome`

Starts the interactive launcher to load tests in `/cypress`.

`$(npm bin)/cypress open`

