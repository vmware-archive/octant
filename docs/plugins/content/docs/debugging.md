---
weight: 90
---

# Debugging

Plugins run as an independent process from the dashboard. A panic within a plugin should not crash a running dashboard.

More detailed logging can be used to debug by passing the verbose flag, `-v`, when running the dashboard.

## Is the plugin registered by the dashboard?

When starting the dashboard, the logs will show a list of registered plugins and their capabilities. If the plugin is not shown as registered in the logs, check the plugin binary is located in the correct plugin path. Make sure the correct GVK is used along with the relevant Capabilities enabled.

```
INFO    plugin/manager.go:286   registered plugin "plugin-name" {"plugin-name": "pluginstub", "cmd": "/home/sfoo/.config/vmdash/plugins/pluginstub", "metadata": {"Name":"plugin-name","Description":"a description","Capabilities":{"SupportsPrinterConfig":[{"Group":"","Version":"v1","Kind":"Pod"}],"SupportsPrinterStatus":[{"Group":"","Version":"v1","Kind":"Pod"}],"SupportsPrinterItems":[{"Group":"","Version":"v1","Kind":"Pod"}],"SupportsObjectStatus":[{"Group":"","Version":"v1","Kind":"Pod"}],"SupportsTab":[{"Group":"","Version":"v1","Kind":"Pod"}]}}}
```

## How to determine if port forward is working?

The UI provides a table of all active port forwarding with links to the running pod. Once a port forward is active, a URL will be available next to the container port.
