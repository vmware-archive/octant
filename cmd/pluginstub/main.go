package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/plugin/api"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type stub struct {
	apiClient api.Service
	mu        sync.Mutex
}

var _ plugin.Service = (*stub)(nil)

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

func (s *stub) Print(object runtime.Object) (plugin.PrintResponse, error) {
	if object == nil {
		return plugin.PrintResponse{}, errors.Errorf("object is nil")
	}

	ctx := context.Background()
	key, err := objectstoreutil.KeyFromObject(object)
	if err != nil {
		return plugin.PrintResponse{}, err
	}
	u, err := s.apiClient.Get(ctx, key)
	if err != nil {
		return plugin.PrintResponse{}, err
	}

	log.Printf("loaded object from objectstore: %v", u)

	msg := fmt.Sprintf("update from plugin at %s", time.Now().Format(time.RFC3339))

	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Status: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
		Items: []component.FlexLayoutItem{
			{
				Width: component.WidthHalf,
				View:  component.NewText("item 1 from plugin"),
			},
			{
				Width: component.WidthFull,
				View:  component.NewText("item 2 from plugin"),
			},
		},
	}, nil
}

func (s *stub) ObjectStatus(object runtime.Object) (plugin.ObjectStatusResponse, error) {
	if object == nil {
		return plugin.ObjectStatusResponse{}, errors.New("object is nil")
	}

	status := component.PodSummary{
		Status:  component.NodeStatusOK,
		Details: []component.Component{component.NewText("status details from plugin")},
	}

	return plugin.ObjectStatusResponse{
		ObjectStatus: status,
	}, nil
}

func (s *stub) PrintTab(object runtime.Object) (*component.Tab, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	layout := flexlayout.New()
	section := layout.AddSection()
	err := section.Add(component.NewText("content from a plugin"), component.WidthHalf)
	if err != nil {
		return nil, err
	}

	tab := component.Tab{
		Name:     "PluginStub",
		Contents: *layout.ToComponent("Plugin"),
	}

	return &tab, nil
}

func main() {
	log.Print("plugin is starting")
	plugin.Serve(&stub{})
}
