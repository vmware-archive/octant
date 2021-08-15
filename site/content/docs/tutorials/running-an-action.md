## Running an Action

An Action is a fundamental concept in Octant. It is the mechanism allowing Octant to run a command to change the
state of the UI or perform some operation against a cluster. Since most Actions are executed upon user interaction
with a component, the need to understand an Action is generally only needed by developers and plugin authors.

### Anatomy of an Action

Octant uses websockets to communicate between server and client. An Action is a websocket message that is sent from
the client to the server. The server expects the message to contain a `type` to specify the Action and a `payload`
containing data relevant to that particular Action.

```json
{
  "type":"action.octant.dev/setContentPath",
  "payload":{
    "contentPath":"overview/namespace/default",
    "params":{}
  }
}
```

The above websocket message is an example of how Octant uses Actions to serve content for a given page. When a user
navigates to a URL path, it is sent from the client, electron or web browser, where the server will send a response
based on the requested content path.

```json
{
  "type":"event.octant.dev/content",
  "data":{},
  "Err":null
}
```

Then the returned message provides the data to generate the content and provide errors, if any.

### Actions in Practice

Most features seen in Octant use Actions. A number of these are not exposed publicly as they perform some operation
against the cluster. Examples of such actions can be found here: [https://pkg.go.
dev/github.
com/vmware-tanzu/octant/internal/octant#pkg-constants](https://pkg.go.dev/github.com/vmware-tanzu/octant/internal/octant#pkg-constants)

Similarly, public [Actions](https://pkg.go.dev/github.com/vmware-tanzu/octant/pkg/action#pkg-constants) are
available for consumption and are also used in Octant core.

Many components prompting user interaction in Octant allow binding to an action. For example, a button is used to
send some data to the server.

```go
component.NewButton("button",
	action.CreatePayload("action.octant.dev/myAction", map[string]interface{}{
	    "hello": "world"
	})
)
```

It is the developer's responsibility to ensure the payload schema is correct in addition to error handling.

### Custom Actions in Plugins

If Octant does not provide an Action doing what you want, plugins can create a new action to run some arbitrary code.

The [sample plugin](https://github.com/vmware-tanzu/octant/blob/192f70c65d78798207896df5331a16d057cc9cdf/cmd/octant-sample-plugin/main.go#L27) provides an example of this by showing the websocket client ID in an application level
alert when clicking a button.

This is a powerful feature to integrate Octant with other tools in the ecosystem.
