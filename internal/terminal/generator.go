/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package terminal

import (
	"context"
	"sort"

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

	details := []component.TerminalDetails{}

	for _, t := range tm.List() {
		if t.Key() == key {
			details = append(details, component.TerminalDetails{
				Container: t.Container(),
				Command:   t.Command(),
				UUID:      t.ID(),
				CreatedAt: t.CreatedAt(),
			})
		}
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].CreatedAt.After(details[j].CreatedAt)
	})

	terminal := component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  component.NewTerminal(key.Namespace, key.Name, details),
		},
	}
	terminals.AddSections(terminal)
	return terminals, nil
}
