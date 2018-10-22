# developer-dash

Kubernetes dashboard for developers

## Running

`$ hcli dash`

## Developing

### Prerequisites

* Go 1.11
* npm
* yarn

### Environment variables

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
