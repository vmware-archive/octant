# Getting Started

## Environment Variables

Octant is configurable through environment variables defined at runtime here are some of the notable variables:

* `KUBECONFIG` - set to non-empty location if you want to set KUBECONFIG with an environment variable.
* `OCTANT_DISABLE_OPEN_BROWSER` - set to a true value if you do not want the browser launched when the dashboard starts up.
* `OCTANT_DISABLE_CLUSTER_OVERVIEW` - set to a true value if you do not want the cluster overview module to load.
* `OCTANT_LISTENER_ADDR` - set to address you want dashboard service to start on. (e.g. `localhost:8080`)
* `OCTANT_ACCEPTED_HOSTS` - set to comma-separated string of hosts to be accepted. (e.g. `demo.octant.example.com,awesome.octant.zr`)
* `OCTANT_LOCAL_CONTENT` - set to a directory and dash will serve content responses from here. An example directory lives in `examples/content`
* `OCTANT_PLUGIN_PATH` - add a plugin directory or multiple directories separated by `:`. Plugins will load by default from `$HOME/.config/octant/plugins`

**Notice:** If using [fish shell](https://fishshell.com), tilde expansion may not occur when using `env` to set environment variables.

### Flags as Variables

All command-line flags can also be passed as environment variables by using all UPPERCASE, replacing the `-` with `_` and prefixing them with `OCTANT_`.
When using CLI flags that enable/disable a feature the following values are considered true and false:

  * **True** - "1", "t", "T", "true", "TRUE", "True"
  * **False** - "0", "f", "F", "false", "FALSE", "False"

Some examples:

 * `--namespace=default` becomes `OCTANT_NAMESPACE=default`
 * `--enable-opencensus` becomes `OCTANT_ENABLE_OPENCENSUS=1`
 * `--disable-cluster-overview` becomes `OCTANT_DISABLE_CLUSTER_OVERVIEW=true`

## Command Line Flags

Octant is configurable through command line flags set at runtime. You can see all of the available options by
running `octant --help`.

```sh
Flags:
      --accepted-hosts string         accepted hosts list
      --client-burst int              maximum burst for client throttle (default 400)
      --client-qps float32            maximum QPS for client (default 200)
      --context string                initial context
      --disable-cluster-overview      disable cluster overview
      --disable-open-browser          disable automatic launching of the browser
      --enable-feature-applications   enable applications feature
  -c, --enable-opencensus             enable open census
  -h, --help                          help for octant
      --klog-verbosity int            klog verbosity level
      --kubeconfig string             absolute path to kubeConfig file (default "/home/wwitzel3/.kube/kind-config-octant")
      --listener-addr string          listener address for the octant frontend
      --local-content string          local content path
  -n, --namespace string              initial namespace
      --plugin-path string            plugin path
      --proxy-frontend string         url to send frontend request to, useful for development
      --ui-url string                 dashboard url
  -v, --verbosity count               verbosity level
```

The verbosity has a special type that is used to parse the flag, which means it can be provided
shorthand by just adding more `v` to equal the level count or with an explicit equal sign.

```sh
-v[vv], --verbosity=count      verbosity level
```

For example

```sh
octant -vvv
```

Is equal to

```sh
octant --verbosity=3
```

## Setting Up a Development Environment

* [Go 1.13 or above](https://golang.org/dl/)
* [node 10.15.0 or above](https://nodejs.org/en/)
* [npm 6.4.1 or above](https://www.npmjs.com/get-npm)
* [rice](https://github.com/GeertJohan/go.rice) - packaging web assets into a binary
* [mockgen](https://github.com/golang/mock) - generating go files used for testing
* [protoc](https://github.com/golang/protobuf) - generate go code compatible with gRPC

These build tools can be installed via Makefile with `go run build.go go-install`.

A development binary can be built by `go run build.go build`.

For UI changes, see the [README]({{ site.gh_repo }}/tree/master/web) located in `web/`.

If Docker and [Drone](/docs/drone) are installed, tests and build steps can run in a containerized environment.

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
