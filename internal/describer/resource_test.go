/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func TestBreadcrumbWithNamespace(t *testing.T) {
	workloadsRoot := ResourceLink{Title: "Workloads", Url: "/overview/namespace/($NAMESPACE)/workloads"}
	expected := []component.TitleComponent{
		component.NewLink("", "Workloads", "/overview/namespace/default/workloads"),
		component.NewLink("", "title", ""),
	}

	breadcrumb := getBreadcrumb(workloadsRoot, "title", "", "default")
	assert.NotNil(t, breadcrumb)
	assert.Len(t, breadcrumb, 2)
	assert.Equal(t, breadcrumb, expected)
}

func TestBreadcrumbWithUrlNamespace(t *testing.T) {
	workloadsRoot := ResourceLink{Title: "Workloads", Url: "/overview/namespace/($NAMESPACE)/workloads"}
	expected := []component.TitleComponent{
		component.NewLink("", "Workloads", "/overview/namespace/default/workloads"),
		component.NewLink("", "title", "/title"),
	}

	breadcrumb := getBreadcrumb(workloadsRoot, "title", "/title", "default")
	assert.NotNil(t, breadcrumb)
	assert.Len(t, breadcrumb, 2)
	assert.Equal(t, breadcrumb, expected)
}

func TestBreadcrumbNoNamespace(t *testing.T) {
	root := ResourceLink{Title: "Cluster Overview", Url: "/cluster-overview"}
	expected := []component.TitleComponent{
		component.NewLink("", "Cluster Overview", "/cluster-overview"),
		component.NewLink("", "title", ""),
	}

	breadcrumb := getBreadcrumb(root, "title", "", "")
	assert.NotNil(t, breadcrumb)
	assert.Len(t, breadcrumb, 2)
	assert.Equal(t, breadcrumb, expected)
}
