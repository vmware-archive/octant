package overview

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/portforward"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
)

// PortForwardListDescriber describes a list of port-forwards
type PortForwardListDescriber struct {
}

// Describe describes a list of port forwards as content
func (d *PortForwardListDescriber) Describe(ctx context.Context, prefix, namespace string, clusterClient cluster.ClientInterface, options DescriberOptions) (component.ContentResponse, error) {
	if options.PortForwardSvc == nil {
		return component.ContentResponse{}, errors.New("portforward service is nil")
	}

	list := component.NewList("Port Forwards", nil)

	tblCols := component.NewTableCols("Name", "Ports", "Age")
	tbl := component.NewTable("Port Forwards", tblCols)
	list.Add(tbl)

	for _, pf := range options.PortForwardSvc.List() {
		t := &pf.Target
		apiVersion, kind := t.GVK.ToAPIVersionAndKind()
		pfRow := component.TableRow{
			"Name":  link.ForGVK(t.Namespace, apiVersion, kind, t.Name, t.Name),
			"Ports": describePortForwardPorts(pf),
			"Age":   component.NewTimestamp(pf.CreatedAt),
		}
		tbl.Add(pfRow)
	}

	return component.ContentResponse{
		ViewComponents: []component.ViewComponent{list},
	}, nil
}

func (d *PortForwardListDescriber) PathFilters() []pathFilter {
	filter := newPathFilter("/portforward", d)
	return []pathFilter{*filter}
}

func NewPortForwardListDescriber() *PortForwardListDescriber {
	return &PortForwardListDescriber{}
}

func describePortForwardPorts(pf portforward.PortForwardState) component.ViewComponent {
	lst := component.NewList("", nil)

	for _, p := range pf.Ports {
		portStr := fmt.Sprintf("%d -> %d", p.Local, p.Remote)
		item := component.NewPortForwardDeleter(
			portStr,
			pf.ID,
			component.NewPortForwardPorts(p.Local, p.Remote),
		)
		lst.Add(item)
	}
	return lst
}
