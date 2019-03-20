package main

import (
	"fmt"
	"time"

	"github.com/heptio/developer-dash/pkg/plugin"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type stub struct{}

var _ plugin.Service = (*stub)(nil)

func (s *stub) Register() (plugin.Metadata, error) {
	podGVK := schema.GroupVersionKind{Version: "v1", Kind: "Pod"}

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

	msg := fmt.Sprintf("update from plugin at %s", time.Now())

	return plugin.PrintResponse{
		Config: []component.SummarySection{
			{Header: "from-plugin", Content: component.NewText(msg)},
		},
	}, nil
}

func main() {
	plugin.Serve(&stub{})
}
