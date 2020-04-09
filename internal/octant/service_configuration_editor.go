/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// ServiceConfigurationEditor edits editors.
type ServiceConfigurationEditor struct {
	store store.Store
}

var _ action.Dispatcher = (*ServiceConfigurationEditor)(nil)

// NewServiceConfigurationEditor creates an instance of ServiceConfigurationEditor.
func NewServiceConfigurationEditor(objectStore store.Store) *ServiceConfigurationEditor {
	editor := &ServiceConfigurationEditor{store: objectStore}
	return editor
}

// ActionName returns the name of this action.
func (s *ServiceConfigurationEditor) ActionName() string {
	return "action.octant.dev/serviceEditor"
}

// Handle edits a service: Supported edits:
//   * selector
func (s *ServiceConfigurationEditor) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", s.ActionName())
	logger.
		With("payload", payload).
		Debugf("received action payload")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	name := key.Name

	selectorsList, err := payload.StringSlice("selectors")
	if err != nil {
		return err
	}

	selector := make(map[string]string)
	for i := range selectorsList {
		parts := strings.SplitN(selectorsList[i], ":", 2)
		if len(parts) != 2 {
			return errors.Errorf("invalid selector %s", selectorsList[i])
		}
		selector[parts[0]] = parts[1]
	}

	fn := func(object *unstructured.Unstructured) error {
		return unstructured.SetNestedStringMap(object.Object, selector, "spec", "selector")
	}

	alertType := action.AlertTypeInfo
	message := fmt.Sprintf("Updated Service %q", name)
	if err := s.store.Update(ctx, key, fn); err != nil {
		alertType = action.AlertTypeWarning
		message = fmt.Sprintf("Unable to update Service %q: %s", name, err)
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)

	return nil
}
