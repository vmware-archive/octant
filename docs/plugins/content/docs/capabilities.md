---
weight: 50
---

# Capabilities

The plugin capabilities are defined in the metadata then the plugin returns the view components to be added to the dashboard.

Using the plugin will rely on two methods: `Print` and `PrintTab`.

## Register

Plugins are registered with their name, description, and capability. The plugin's API address is determined by the dashboard server once the plugin begins to run.

```go
func (s *stub) Register(apiAddress string) (plugin.Metadata, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	client, err := api.NewClient(apiAddress)
	if err != nil {
		return plugin.Metadata{}, errors.Wrap(err, "create api client")
	}

	s.apiClient = client

	podGVK := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}
	log.Println("the dashboard plugin api is at", apiAddress)

	return plugin.Metadata{
		Name:        "plugin-name",
		Description: "a description",
		Capabilities: plugin.Capabilities{
			SupportsPrinterConfig: []schema.GroupVersionKind{podGVK},
			SupportsPrinterStatus: []schema.GroupVersionKind{podGVK},
			SupportsPrinterItems:  []schema.GroupVersionKind{podGVK},
			SupportsObjectStatus:  []schema.GroupVersionKind{podGVK},
			SupportsTab:           []schema.GroupVersionKind{podGVK},
			IsModule:              false,
		},
	}, nil
}
```

Plugins can also act as Octant modules. To enable module support, set `IsModule` to true. When acting as modules, plugins can provide
content and navigation entries.

## Summary

A Summary consists of a Config, Status, and Items. The Content can be any of the various components found in [reference](/docs/reference).

```go
func (s *stub) Print(object runtime.Object) (*component.Tab, error) {
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

## Tabs

Adding a new tab via plugin requires a new flexlayout then Tab component. The Name is used in the URL query param, and Contents defines the tab name within the dashboard.

```go
func (s *stub) PrintTab(object runtime.Object) (*component.Tab, error) {
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

## Port Forwarding

A port forward requires a request containing the namespace, pod name, and port number.

```go
	pfRequest := api.PortForwardRequest{
		Namespace: "heptio-contour",
		PodName:   name,
		Port:      6060,
	}
	pfResponse, err := p.apiClient.PortForward(ctx, pfRequest)
	if err != nil {
		return nil, err
	}

	defer p.apiClient.CancelPortForward(ctx, pfResponse.ID)
```

## Navigation

If a plugin is configured as a module, it can supply navigation entries. These navigation entries will displayed with the application's
navigation.

```go
func (p *modulePlugin) Navigation() (navigation.Navigation, error) {
	return navigation.Navigation{
		Title: "Module Plugin",
		Path:  path.Join("/content", "module-plugin", "/"),
		Children: []navigation.Navigation{
			{
				Title:    "Nested Once",
				Path:     path.Join("/content", "module-plugin", "nested-once"),
				IconName: "folder",
				Children: []navigation.Navigation{
					{
						Title:    "Nested Twice",
						Path:     path.Join("/content", "module-plugin", "nested-once", "nested-twice"),
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

If a plugin is configured as a module, it can serve content. The content consists of Octant components. wrapped in a `ContentResponse`.
The function will receive the currently requested content path and can display content based on that path. 

```go
func (p *modulePlugin) Content(contentPath string) (component.ContentResponse, error) {
	return component.ContentResponse{
		Components: []component.Component{
			component.NewText(fmt.Sprintf("hello from plugin: path %s", contentPath)),
		},
	}, nil
}
```
