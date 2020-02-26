# Combined Logging

Support the ability to view all of the container logs for a given workload. Currently you must request each container
log separately.

## Goals
 - Combine the logging output from multiple containers.
 - Tag/mark which log lines are from which container.
 - Provide a list of which container logs are included in the stream.
 - Only send content when there are new log entries.

## Non-Goals
 - Server-side filtering

## Log Streaming

### Current implementation
Our existing log "streaming" is done via a log endpoint with a JSON request. We explicitly register a container logs
handler in `internal/api/api.go`.

    s.HandleFunc("/logs/namespace/{namespace}/pod/{pod}/container/{container}", containerLogsHandler(ctx, a.dashConfig))

When visiting the Logs tab we generate an XHR matching the URL pattern we've registered for container logs. This XHR
request is generated using the selected container from the listed of containers. This list of containers is provided
by Octant to the Log component. A user can only view a single containers logs at a time.

Example request:

    http://localhost:7777/api/v1/logs/namespace/octant-demo/pod/kuard-7bbd76c779-6n2p4/container/kuard

Octant responds with a JSON object that contains the entirety of the log data each request. The request is repeated via
a polling mechanism in the PodLogsStreamer, every 5 second. When a new container is selected from the dropdown, the
existing poll is canceled and a new one is created with a request that uses the newly selected container.

### Suggested implementation
Add a sentinel value of `AllContainers` to Octant and make this the default choice when requesting logs.

Add a set of new actions to Octant to handle starting a log stream.

Add a new EventType to Octant for log streams that the WebsocketService can subscribe to.

Extend our existing log entry to include the container the log entry is from.

#### Stream Aggregation
Create a go routine that wraps the `Stream` call for each containers Octant will stream logs for.

    client.CoreV1().Pods(lp.namespace).GetLogs(lp.podName, &corev1.PodLogOptions{
    		Container:  lp.container,
    		Follow:     false,
    		Timestamps: true,
    	}).Stream()
    }

The `Stream` call returns an `io.ReadCloser`. We cannot use the built-in `io.MultiReader` to produce a reader from all
of our log streams because we need to identify which container each log line is from.

Create a multi-reader like function that takes a map of `containerId:Reader` and a channel to send log entries. Our
multi-reader function will iterate over the map of readers, create scanners, and produce log entry objects
and send them to the output channel.

A worker loop will read the log entries from the output channel and write them to the websocket client using the
new EventType for log streams.

Ensure the worker loop respects the given Context and properly cancels streaming when navigating away from the logging
tab and also ensure we are calling `stream.Close` using a `defer`.

#### Frontend
Make our log component on the frontend append data from the WebsocketService instead of redrawing
the component every time. Similar to how we send data for the Terminal component.

Change the frontend component so that the container selector defaults to our all containers value and that selecting
a container filters log lines by using the container in the log entry config.
