/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// ContainerEditor edits containers.
type ContainerEditor struct {
	store store.Store
}

var _ action.Dispatcher = (*ContainerEditor)(nil)

// NewContainerEditor creates an instance of ContainerEditor.
func NewContainerEditor(objectStore store.Store) *ContainerEditor {
	editor := &ContainerEditor{
		store: objectStore,
	}

	return editor
}

// ActionName returns name of this action.
func (e *ContainerEditor) ActionName() string {
	return "overview/containerEditor"
}

// Handle edits a container. Supported edits:
//   * image
func (e *ContainerEditor) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := internalLog.From(ctx).With("actionName", e.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	containersPathData, err := payload.String("containersPath")
	if err != nil {
		return err
	}

	var containersPath []string
	if err := json.Unmarshal([]byte(containersPathData), &containersPath); err != nil {
		return err
	}

	containerName, err := payload.String("containerName")
	if err != nil {
		return err
	}

	containerImage, err := payload.String("containerImage")
	if err != nil {
		return err
	}

	fn := updateContainer(containersPath, logger, containerName, containerImage)

	message := fmt.Sprintf("Container %q was updated", containerName)
	alertType := action.AlertTypeInfo
	if err := e.store.Update(ctx, key, fn); err != nil {
		message = fmt.Sprintf("Unable to update container %q: %s", containerName, err)
		alertType = action.AlertTypeWarning
		logger := internalLog.From(ctx)
		logger.WithErr(err).Errorf("update container")
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)

	alerter.SendAlert(alert)
	return nil
}

func updateContainer(containersPath []string, logger log.Logger, containerName string, containerImage string) func(object *unstructured.Unstructured) error {
	return func(object *unstructured.Unstructured) error {
		containersRaw, found, err := unstructured.NestedSlice(object.Object, containersPath...)
		if err != nil {
			return err
		}

		if !found {
			logger.Warnf("unable to find containers within object")
			return nil
		}

		var updatedContainers []interface{}

		for _, containerRaw := range containersRaw {
			container, ok := containerRaw.(map[string]interface{})
			if !ok {
				return errors.New("unable to parse containers format")
			}
			name, found, err := unstructured.NestedString(container, "name")
			if err != nil {
				return errors.Wrap(err, "looking for container name")
			}

			if found && name == containerName {
				container["image"] = containerImage
			}

			updatedContainers = append(updatedContainers, container)
		}

		return unstructured.SetNestedSlice(object.Object, updatedContainers, containersPath...)
	}
}
