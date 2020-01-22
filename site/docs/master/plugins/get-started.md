# Get Started

## Build an example plugin

An example plugin can be found within the developer dashboard repo.

Install the plugin using:

```sh
go run build.go install-test-plugin
```

Alternatively, build the go binary using `go build` then move the binary to the install path described below.

## Installation

'go run build.go install-test-plugin' installs the plugin by creating a `$HOME/.config/octant/plugin/` directory then building the binary to that location.

Run plugins from additional paths by setting paths to the `OCTANT_PLUGIN_PATH` environment variable when running the dashboard.

Octant will also respect `XDG_CONFIG_HOME` on Unix and `LocalAppData` on Windows for default plugin paths.

## Define Capability

Each plugin must have a defined name, description, and capability.

<!-- TODO: naming restrictions or conventions -->

Plugins can provide a `PrintResponse` containing capabilities enabled by a provided GVK.

### Config

A plugin with support for `PrinterConfig` appends a view component to the Configuration table of the supported GVK(s).

The header is added to the column on the left. Content is a component that is added to the right.

![PrinterConfig](kuard_deployment_config.png)

Certain GVK such as Deployments have a Configuration but not Status.

### Status

A plugin with support for `PrinterStatus` appends a view component to the Status table of the supported GVK(s).

![PrinterStatus](kuard_pod_config_status.png)

This pod has both a Configuration and Status.

### Items

A plugin with support for `PrinterItems` allow adding a `FlexLayoutItem` consisting of a width and a view component.

## Uninstall

Plugins can be removed by deleting the plugin binary from `~/.config/octant/plugins`. An example of deleting a plugin is shown below 
where `octant-sample-plugin` is the plugin that will be uninstalled:

```
rm ~/.config/octant/plugins/octant-sample-plugin
```

After deleting the plugin binary and restarting Octant, you should no longer see the plugin available as part of the Octant dashboard. 
