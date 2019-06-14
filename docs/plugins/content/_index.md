---
title: Overview
type: docs
---

# Introduction to Plugins

Plugins are binaries that run alongside developer dashboard to provide additional functionality. Plugins are built using [go-plugin](https://github.com/hashicorp/go-plugin) in order to communicate with the dashboard over gRPC.

The goal of this documentation is to provide:

1. Instructions for building and installing a plugin
2. A description of plugin capabilities
3. Code examples for developers to write their own plugins

## What can plugins do?

Plugins can:

 * Add new tabs to the dashboard
 * Include additional content to an existing summary section
 * Create a new section in an existing tab
 * Port forward to a running pod

By using the components available in `/pkg/view/components`, plugins can display new information in desired areas of the dashboard. Entirely new layouts can also be created through `pkg/view/flexlayout`.

## Next Steps

List channels for feedback/troubleshooting.

 * Slack channel (#dev-dashboard)
 * Email: <project-devdash@groups.vmware.com>

