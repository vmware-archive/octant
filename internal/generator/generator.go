/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package generator

import (
	"context"

	kLabels "k8s.io/apimachinery/pkg/labels"

	"github.com/vmware/octant/internal/portforward"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/view/component"
)

type Options struct {
	LabelSet       *kLabels.Set
	PortForwardSvc portforward.PortForwarder
	PluginManager  *plugin.Manager
}

type Generator interface {
	Generate(ctx context.Context, path, prefix, namespace string, opts Options) (component.ContentResponse, error)
}
