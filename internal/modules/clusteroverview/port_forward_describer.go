/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"
	"fmt"

	"github.com/vmware/octant/internal/describer"
	"github.com/vmware/octant/internal/portforward"
	"github.com/vmware/octant/pkg/view/component"
)

// PortForwardListDescriber describes a list of port-forwards
type PortForwardListDescriber struct {
}

func NewPortForwardListDescriber() *PortForwardListDescriber {
	return &PortForwardListDescriber{}
}

// Describe describes a list of port forwards as content
func (d *PortForwardListDescriber) Describe(ctx context.Context, prefix, namespace string, options describer.Options) (component.ContentResponse, error) {
	portForwarder := options.PortForwarder()

	list := component.NewList("Port Forwards", nil)

	tblCols := component.NewTableCols("Name", "Ports", "Age")
	tbl := component.NewTable("Port Forwards", tblCols)
	list.Add(tbl)

	for _, pf := range portForwarder.List() {
		t := &pf.Target
		apiVersion, kind := t.GVK.ToAPIVersionAndKind()
		nameLink ,err := options.Link.ForGVK(t.Namespace, apiVersion, kind, t.Name, t.Name)
		if err != nil {
			return describer.EmptyContentResponse, err
		}

		pfRow := component.TableRow{
			"Name":  nameLink,
			"Ports": describePortForwardPorts(pf),
			"Age":   component.NewTimestamp(pf.CreatedAt),
		}
		tbl.Add(pfRow)
	}

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func (d *PortForwardListDescriber) PathFilters() []describer.PathFilter {
	filter := describer.NewPathFilter("/port-forward", d)
	return []describer.PathFilter{*filter}
}

func describePortForwardPorts(pf portforward.State) component.Component {
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
