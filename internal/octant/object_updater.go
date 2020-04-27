/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

package octant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// ObjectUpdateFromPayload loads an object from the payload.
// The object source in YAML format should exist in the `update` key.
func ObjectUpdateFromPayload(payload action.Payload) (*unstructured.Unstructured, error) {
	s, err := payload.String("update")
	if err != nil {
		return nil, fmt.Errorf("read object source from payload: %w", err)
	}

	object, err := kubernetes.ReadObject(strings.NewReader(s))
	if err != nil {
		return nil, fmt.Errorf("read object from payload: %w", err)
	}

	return object, nil
}

type ObjectUpdaterDispatcherOption func(dispatcher *ObjectUpdaterDispatcher)

// ObjectUpdaterDispatcher is an action that updates an object.
type ObjectUpdaterDispatcher struct {
	store             store.Store
	objectFromPayload func(payload action.Payload) (*unstructured.Unstructured, error)
}

var _ action.Dispatcher = &ObjectUpdaterDispatcher{}

// NewObjectUpdaterDispatcher creates an instance of ObjectUpdaterDispatcher.
func NewObjectUpdaterDispatcher(objectStore store.Store, options ...ObjectUpdaterDispatcherOption) *ObjectUpdaterDispatcher {
	o := ObjectUpdaterDispatcher{
		store:             objectStore,
		objectFromPayload: ObjectUpdateFromPayload,
	}

	for _, option := range options {
		option(&o)
	}

	return &o
}

// ActionName returns the action name this dispatcher responds to.
func (o ObjectUpdaterDispatcher) ActionName() string {
	return ActionUpdateObject
}

// Handle updates an object using a payload if possible.
func (o ObjectUpdaterDispatcher) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx)
	expiration := time.Now().Add(10 * time.Second)

	object, err := o.objectFromPayload(payload)
	if err != nil {
		sendAlert(
			alerter,
			action.AlertTypeError,
			fmt.Sprintf("load object from payload: %v", err.Error()),
			&expiration)
		return nil
	}

	key, _ := store.KeyFromPayload(payload)
	err = o.store.Update(ctx, key, func(u *unstructured.Unstructured) error {
		if object.GetAPIVersion() != u.GetAPIVersion() {
			return fmt.Errorf("object API version cannot be updated")
		}
		if object.GetKind() != u.GetKind() {
			return fmt.Errorf("object kind cannot be updated")
		}
		if object.GetName() != u.GetName() {
			return fmt.Errorf("object name cannot be updated")
		}

		delete(object.Object, "status")

		for k := range object.Object {
			u.Object[k] = object.Object[k]
		}
		return nil
	})

	if err != nil {
		sendAlert(
			alerter,
			action.AlertTypeError,
			fmt.Sprintf("update object: %s", err.Error()),
			&expiration)

		logger.WithErr(err).Errorf("update object")
		return nil
	}

	successMessage := fmt.Sprintf("Updated %s (%s) %s in %s",
		object.GetKind(),
		object.GetAPIVersion(),
		object.GetName(),
		object.GetNamespace())
	sendAlert(alerter, action.AlertTypeInfo, successMessage, &expiration)

	return nil
}
