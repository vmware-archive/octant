# developer-dash

Kubernetes dashboard for developers

## Running

Start the developer dashboard:

`$ sugarloaf dash`

Check the version:

`$ sugarloaf version`

### Prerequisites

* Go 1.11
* npm 6.4.1 or higher
* yarn
* [rice CLI](https://github.com/GeertJohan/go.rice)
  * Install with `go get github.com/GeertJohan/go.rice/rice`
* [mockgen](https://github.com/golang/gomock)
  * go get github.com/golang/mock/gomock && go install github.com/golang/mock/mockgen

## Install

### Download a prebuilt binary

Go to the [releases page](https://github.com/heptio/developer-dash/releases) and download the tarball.

Extract the tarball:

```sh
$ tar -xzvf ~/Downloads/sugarloaf_0.3.0_Linux-64bit.tar.gz
sugarloaf_0.3.0_Linux-64bit/README.md
sugarloaf_0.3.0_Linux-64bit/sugarloaf
```

Verify it runs:

`$ ./sugarloaf_0.3.0_Linux-64bit/sugarloaf version`

Decide to move the binary in `/usr/local/bin` or your home directory. Installing to `/usr/local/bin` is for system-wide installation but makes running multiple versions difficult. If the dashboard is installed to your home directory, make sure to update your `$PATH` variable then check `which sugarloaf` to verify installation is successful.

### Manually build and install

This option is for users who want to build from master. Make sure the prerequisites listed above are installed.

`$ go get github.com/heptio/developer-dashboard`

Package the web assets to be built into the binary.

`$ make web-build`

There should be a new directory: `$GOPATH/src/github.com/heptio/developer-dash/web/build`. Finally, build the binary:

`$ make sugarloaf-dev`

The `sugarloaf` binary will be found in `$GOPATH/src/github.com/heptio/developer-dash/build`.

### Environment variables

* `KUBECONFIG` - set to non-empty location if you want to set KUBECONFIG with an environment variable.
* `DASH_DISABLE_OPEN_BROWSER` - set to a non-empty value if you don't the browser launched when the dashboard start up.
* `DASH_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)
* `DASH_VERBOSE_CACHE` - set to a non-empty value to view cache actions
* `DASH_LOCAL_CONTENT` - set to a directory and dash will serve content responses from here. An example directory lives in `examples/content`

### Running development web UI

`$ make setup-web`

### Running development server

The development server allows running the dashboard while monitoring changes in `/web`.

Start the dashboard running on a development server:

`$ make -j ui-client ui-server`

Navigate to `localhost:7777` on a browser to view cluster data.

### Caveats

* If using [fish shell](https://fishshell.com), tilde expansion may not occur when using `env` to set environment variables.
