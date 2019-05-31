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
		},
	}, nil
}
```

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
