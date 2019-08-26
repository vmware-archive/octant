/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package clusteroverview

import (
	"context"

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

var _ describer.Describer = (*PortForwardListDescriber)(nil)

// Describe describes a list of port forwards as content
func (d *PortForwardListDescriber) Describe(ctx context.Context, prefix, namespace string, options describer.Options) (component.ContentResponse, error) {
	portForwarder := options.PortForwarder()

	list := component.NewList("Port Forwards", nil)

	tblCols := component.NewTableCols("Name", "Namespace", "Ports", "Age")
	tbl := component.NewTable("Port Forwards", "There are no port forwards!", tblCols)
	list.Add(tbl)

	for _, pf := range portForwarder.List(ctx) {
		t := &pf.Target
		apiVersion, kind := t.GVK.ToAPIVersionAndKind()
		nameLink, err := options.Link.ForGVK(t.Namespace, apiVersion, kind, t.Name, t.Name)
		if err != nil {
			return describer.EmptyContentResponse, err
		}

		pfRow := component.TableRow{
			"Name":      nameLink,
			"Namespace": component.NewText(t.Namespace),
			"Ports":     component.NewPorts(describePortForwardPorts(pf)),
			"Age":       component.NewTimestamp(pf.CreatedAt),
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

func (d *PortForwardListDescriber) Reset(ctx context.Context) error {
	return nil
}

func describePortForwardPorts(pf portforward.State) []component.Port {
	var list []component.Port
	apiVersion, kind := pf.Target.GVK.ToAPIVersionAndKind()
	pfs := component.PortForwardState{}

	for _, p := range pf.Ports {
		pfs.ID = pf.ID
		pfs.Port = int(p.Local)
		pfs.IsForwarded = true

		port := component.NewPort(
			pf.Target.Namespace,
			apiVersion,
			kind,
			pf.Target.Name,
			int(p.Remote),
			"TCP", pfs)
		list = append(list, *port)
	}
	return list
}
