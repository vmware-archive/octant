## Overview

Octant is a tool made to enhance developer experience on Kubernetes and allow a new approach for understanding
complex, distributed systems. The goal of these docs is to explain key concepts that will help get the most value
out of using Octant. A breadth of features are available whether the user is someone new to Kubernetes or an
experienced engineer debugging issues in production.

This site is a constant work in progress capturing an actively developed project. If something is missing or have
room for improvement, feel free to contribute through GitHub, Slack, or attending the community meeting.

## Helpful Links

The tools below are standard for local development with Kubernetes. These are recommended for installation alongside
usage of Octant.

- [Docker](https://github.com/docker) - Containers for development and production
- [KinD](https://github.com/kubernetes-sigs/kind) - Simplified local cluster environment
- [kubectl](https://github.com/kubernetes/kubectl) - CLI for interacting with a Kubernetes cluster
- [Helm](https://github.com/helm/helm) - Install applications on a cluster similar to a package manager

For developers who are comfortable in Go, [client-go](https://github.com/kubernetes/client-go) is a useful library
for programmatically interacting with a cluster and understanding how Octant runs under the hood.

## Getting Started

Once you have downloaded Octant, the best way to start is running Octant to explore how the UI shows an existing
cluster. Additional plugins can enhance the experience as well as introduce potentially new tools for your workflows.

For users new to Kubernetes and Octant, create a local cluster using KinD and run some commands with kubectl along
tutorials found on [https://kubernetes.io/](https://kubernetes.io/). As you become more familiar with various
commands, see if they can be replaced through Octant instead.

Finally, Octant has a plugin model that allows building custom UI through its component library.
