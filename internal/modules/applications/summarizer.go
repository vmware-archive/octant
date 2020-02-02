package applications

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

//go:generate mockgen -destination=./fake/mock_summarizer.go -package=fake github.com/vmware-tanzu/octant/internal/modules/applications Summarizer

var (
	applicationListColumns = component.NewTableCols("Name", "Instance", "Version", "State", "Pods")
)

const (
	appLabelName     = "home.kubernetes.io/name"
	appLabelInstance = "home.kubernetes.io/instance"
	appLabelVersion  = "home.kubernetes.io/version"
)

// SummarizerConfig is configuration for Summarize.
type SummarizerConfig interface {
	ObjectStore() store.Store
}

// Summarizer summarizes applications for a namespace. Applications are a group of objects
// labeled with matching application labels. Application labels are:
//   * home.kubernetes.io/name
//   * home.kubernetes.io/instance
//   * home.kubernetes.io/version
type Summarizer interface {
	// Summarize generates a table summary.
	Summarize(ctx context.Context, namespace string, config SummarizerConfig) (*component.Table, error)
}

type summarizer struct{}

// Summarize converts applications in namespace to a table.
func (s *summarizer) Summarize(ctx context.Context, namespace string, config SummarizerConfig) (*component.Table, error) {
	if config == nil {
		return nil, errors.Errorf("config is nil")
	}

	applications, err := listApplications(ctx, config.ObjectStore(), namespace)
	if err != nil {
		return nil, err
	}

	table := component.NewTable("Applications", "applications", applicationListColumns)
	for _, d := range applications {
		table.Add(component.TableRow{
			"Name":     component.NewLink("", d.Name, d.Path("applications", namespace)),
			"Instance": component.NewText(d.Instance),
			"Version":  component.NewText(d.Version),
			"State":    component.NewText("state"),
			"Pods":     component.NewText(fmt.Sprintf("%d", d.PodCount)),
		})
	}

	return table, nil
}
