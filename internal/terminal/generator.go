/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenerateComponent generates the Terminal components from the given details.
func GenerateComponent(ctx context.Context, tm Manager, object runtime.Object) (component.Component, error) {
	terminals := component.NewFlexLayout("Terminals")

	key, err := store.KeyFromObject(object)
	if err != nil {
		return nil, err
	}

	for _, t := range tm.List(ctx) {
		if t.Key() == key {
			// Create terminal
			terminal := component.FlexLayoutSection{
				{
					Width: component.WidthFull,
					View:  component.NewTerminal(key.Namespace, key.Name, t.Container(), t.Command(), t.ID()),
				},
			}

			// consider extending FlexLayout with AddTabs
			terminals.AddSections(terminal)
		}
	}
	return terminals, nil
}
