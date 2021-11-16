# Octant Container Exec

Add the ability to execute a command in a container.

## Summary

## Motivation

Currently executing a command in a container would be done using a tool such as `kubectl`. Executing a command in a container includes having to pass multiple arguments such as pod name, command, and namespace to `kubectl exec`. Octant intends to make this an easy one or two click operation from the Pod view.

### Goals
 - Allow interactive and non-interactive commands to be executed against a container.

### Non-Goals
 - Aggregating results from multiple containers.
 - Being a generic interactive shell to launch `kubectl` from.
 - Executing a single command across multiple containers at once.

## Proposal

### Implementation Details

#### Container Exec Action

An action that holds the details of the command to execute will be created.

This action will be generated from the Pod view and include the container(s), the command to run, and if the resulting terminal should be interactive interactive.

The result of the action will be rendered as TerminalComponent(s). These will be displayed as new nested tabs on the Pod views Execute tab.

The tabs will be flagged closable and have a close action that will end the context for interactive sessions.

#### TabComponent

Add an `isClosable` and `closeAction` to tabs.

#### TerminalComponent

Create a new TerminalComponent using https://xtermjs.org/ for the frontend.

The backend component will include the stream name for the component which will associate the output from the executor with a given TerminalComponent.

#### Container Command Executor

Create an Executor using https://github.com/kubernetes/client-go/blob/master/tools/remotecommand/remotecommand.go

We will create a factory that assembles `SPDYExecutor` given a desired container and command to execute.

More specifically that factory will use `remotecommand.NewSPDYExecutor` to create a `SPDYExecutor` and create unique stream names for each executor. Calling the `Stream` method we will send output to our `TerminalComponent` using websockets.

These executors will have contexts to ensure they are cleaned up properly when the terminal goes away or Octant is closed.

#### User Interface

Add a listing page of all currently open temrinals. This page will allow you easily jump to an open terminal as well as close open terminals (including disconnected terminals).

Execute tab will be added to Pod resource. Add executes are run, new, nested terminal windows will be opened in this tab.
The terminal tabs will be named: container_name_%(count) where count is the number of exec commands currently open for a given container.

An Execute icon/button will be placed on the container card. Clicking this icon will result in an Execute dialog for the container:

   - clearly state the container the command is being executed against.
   - command to execute input box
   - interactive check box
   - Execute / Cancel buttons


#### Disconnects and Timeouts

Terminals will timeout after a set idle period of 300 seconds. Alerts will be used to notify a user of when a terminal has timed out or disconnected. Terminals that are disconnected due to timeout or communication error will remain present in the UI.

This will also include containers that go away.
