/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateComponent generates the Terminal components from the given details.
func GenerateComponent(ctx context.Context, tm Manager, object runtime.Object) (component.Component, error) {
	terminals := component.NewFlexLayout("Terminals")
	return terminals, nil
}
