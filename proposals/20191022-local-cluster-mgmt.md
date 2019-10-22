# Octant Local Cluster Management

Add the ability to manage local KIND clusters through Octant.

## Summary

## Motivation

Reduce the overhead of managing local clusters. While using the command line to create a local KIND clusters is simple, once those
clusters have been created, managing those clusters becomes a burden, especially after any time away.

### Goals
 - KIND cluster verbs: create, delete, stop, start
 - Listing KIND clusters.
 - Single click to switch Octant  context.

### Non-Goals
 - Managing non-KIND clusters.

### Optional
  - Install KIND from Octant.

## Proposal

### Implementation Details

This feature will depend on the module navigation feature.

Frontend:
  - Create a new module to act as a container for local cluster management.
  - Listing view that supports KIND actions
  - Octant action that switches the current context to a selected cluster.

Backend:
  - Use sigs.k8s.io/kind/pkg/cluster to implement KIND actions on the backend (See priror art)
    - List
    - Create
    - Delete
    - Start
    - Stop
  - Prior Art: https://github.com/bryanl/k8s-lab/blob/master/pkg/k8slab/cluster.go 
