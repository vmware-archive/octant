# developer-dash

Kubernetes dashboard for developers

## Running

`$ hcli dash`

## Developing

### Prerequisites

* Go 1.11
* npm
* yarn
* [rice CLI](https://github.com/GeertJohan/go.rice)
  * Install with `go get github.com/GeertJohan/go.rice/rice`

### Environment variables

* `KUBECONFIG` - set to non-empty location if you want to set KUBECONFIG with an environment variable.

* `DASH_DISABLE_OPEN_BROWSER` - set to a non-empty value if you don't the browser launched when the dashboard start up.
* `DASH_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)

* `DASH_VERBOSE_CACHE` - set to a non-empty value to view cache actions

* `DASH_TELEMETRY_ADDRESS` - set telemetry address (defaults to `telemetry.corp.heptio.net:443`)
* `DASH_DISABLE_TELEMETRY` - set to non-empty value to disable telemetry

### Running development web UI

`$ make setup-web`

### Building binary with embedded web assets

1) Run `$ make web-build` to rebuild web assets

2) Create `./build/hcli`: `$ make hcli-dev`

### Caveats

* If using [fish shell](https://fishshell.com), tilde expansion may not occur when using `env` to set environment variables.