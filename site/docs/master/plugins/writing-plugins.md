# Writing Plugins

When you want to extend Octant to do something that is not part of the core functionality you will need to write a plugin. Writing an Octant plugin consists of three main parts: defining the capabilities, creating handlers, and registering and serving the plugin.

## Capabilities

Using `plugin.Capabilities` you can define your desired list of capabilites using GVKs. Octant provides a set of well defined capabilites for plugins. These capabilites directly map to Octant renderers and allow your plugin to inject its own components in to the view.

When `plugin.Metadata.IsModule` to true plugins can provide content and navigation entries.

```go
capabilities := &plugin.Capabilities{
	SupportsTab:           []schema.GroupVersionKind{podGVK},
	IsModule:              False,
}
```

The above defines a non-module plugin that will generate a new tab for Pod objects.

## Handlers

Using `service.HandlerFuncs` you will assign handler functions for each of the capabilities for your plugin.

```go
func handleTab(dashboardClient service.Dashboard, object runtime.Object) (*component.Tab, error) {
	// ...
}

handlers := service.HandlerFuncs{
	PrintTab: handleTab,
}
```

### Handling Actions
In Octant you can create custom action handlers that you can trigger from button actions in the UI. There are also
built-in actions which are triggered from internal Octant events, those are defined in [octant/pkg/action/action.go]((https://github.com/vmware-tanzu/octant/blob/master/pkg/action/action.go).

Here is an example of setting up your plugin to know when the current namespace has changed.

```go
	capabilities := &plugin.Capabilities{
		ActionNames:           []string{action.RequestSetNamespace},
	}
	// Set up the action handler.
	options := []service.PluginOption{
		service.WithActionHandler(handleAction),
	}

	func handleAction(request *service.ActionRequest) error {
		switch request.ActionName {
			case action.RequestSetNamespace:
				namespace, err := action.Payload.String("namespace")
				// err check, do work
		}
		return nil
	}
```

## Register and Serve

Registering and serving your plugin is the final step to get your plugin communicating with Octant. This is also where you
will pass in the name and description for the plugin.

```go
p, err := service.Register("plugin-name", "a description", capabilities, handlers)
if err != nil {
	log.Fatal(err)
}

log.Printf("octant-sample-plugin is starting")
p.Serve()
```


## Example

Octant ships with an [example plugin](https://github.com/vmware-tanzu/octant/blob/master/cmd/octant-sample-plugin/main.go).

## More About Capabilities

Octant provides a well defined set of capabilites for plugins to implement. These include:

* Print support: printing config, status, and items to the overview summary for an object.
* Tab support: creating a new tab in the overview for an object.
* Object status: adding object status to a given object.
* Actions: defining custom actions that route to the plugin.

For plugins that as configured as modules the capabilities also include:

* Navigation support; adding entries to the navigation section.
* Content support; creating content to display on a given path.

## Print

A `PrintResponse` consists of a Config, Status, and Items. The Content can be any of the various components found in [reference](/docs/reference).

```go
func handlePrint(dashboardClient service.Dashboard, object runtime.Object) (*plugin.PrintResponse, error) {
	...
	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText("")},
		},
		Status: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText("")},
		},
		Items: []component.FlexLayoutItem{
			{
				Width: component.WidthFull,
				View:  component.NewText(""),
			},
		},
	}, nil
```

## Tab

Adding a new tab via a plugin requires a new flexlayout then Tab component. The Name is used in the URL query param, and Contents defines the tab name within the dashboard.

```go
func handleTab(dashboardClient service.Dashboard, object runtime.Object) (*component.Tab, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	layout := flexlayout.New()

	tab := component.Tab{
		Name:     "Plugin",
		Contents: *layout.ToComponent("Plugin Tab Name"),
	}

	return &tab, nil
}
```

## Object Status

An `ObjectStatusResponse` has an `ObjectStatus` which currently maps to a `PodSummary` and contains a list of Details and a NodeStatus (ok, warning, error). Details can be any of the various components found in [reference](/docs/reference).

```go
func handleObjectStatus(dashboardClient service.Dashboard, object runtime.Object) (plugin.ObjectStatusResponse, error) {
	if object == nil {
		return plugin.ObjectStatusResponse{}, errors.New("object is nil")
	}

	objectStatusResp := plugin.ObjectStatusResponse{
		ObjectStatus: component.PodSummary{
			Details: []component.Component{component.NewText("plugin-name: added status")},
			Status:  component.NodeStatusOK,
		},
	}

	return objectStatusResp, nil
}
```

## Actions

## Navigation

Plugins configured as modules can supply navigation entries. These navigation entries will be displayed with the application's
navigation.

```go
var pluginName = "plugin-name"
var pluginPath = path.Join("content", pluginName)

func handleNavigation(dashboardClient service.Dashboard) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title: "Module Plugin",
		Path:  path.Join(pluginPath, "/"),
		Children: []navigation.Navigation{
			{
				Title:    "Nested Once",
				Path:     path.Join(pluginPath, "nested-once"),
				IconName: "folder",
				Children: []navigation.Navigation{
					{
						Title:    "Nested Twice",
						Path:     path.Join(pluginPath, "nested-once", "nested-twice"),
						IconName: "folder",
					},
				},
			},
		},
		IconName: "cloud",
	}, nil

}
```

## Content

Plugins configured as modules can serve content. The content consists of Octant components wrapped in a `ContentResponse`.
The function will receive the currently requested content path and can display content based on that path. 

```go
func handleContent(dashboardClient service.Dashboard, contentPath string) (component.ContentResponse, error) {
	return component.ContentResponse{
		Components: []component.Component{
			component.NewText(fmt.Sprintf("hello from plugin: path %s", contentPath)),
		},
	}, nil
}
```

## Module Path

Currently Octant creates a non-configurable base path for your plugin that is derived from the name of the plugin.

```sh
/content/plugin-name
```

You can create nested paths that route to your module using that base path. Plugins should handle nested paths in the `Content` function and dispatch the responses accordingly.
