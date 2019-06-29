/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package overview

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"github.com/vmware/octant/internal/action"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/store"
)

const (
	configurationEditorAction = "deployment/configuration"
)

type ClusterClient interface {
	Resource(kind schema.GroupKind) (schema.GroupVersionResource, error)
	DynamicClient() (dynamic.Interface, error)
}

type ConfigurationEditor struct {
	logger log.Logger
	store  store.Store
}

func NewConfigurationEditor(logger log.Logger, objectStore store.Store) *ConfigurationEditor {
	return &ConfigurationEditor{
		logger: logger,
		store:  objectStore,
	}
}

func (e *ConfigurationEditor) ActionName() string {
	return configurationEditorAction
}

func (e *ConfigurationEditor) Handle(ctx context.Context, payload action.Payload) error {
	e.logger.With("payload", payload, "actionName", "deployment/configuration").Infof("received action payload")

	gvk, err := payload.GroupVersionKind()
	if err != nil {
		return err
	}

	name, err := payload.String("name")
	if err != nil {
		return err
	}

	namespace, err := payload.String("namespace")
	if err != nil {
		return err
	}

	replicaCountFloat, err := payload.Float64("replicas")
	if err != nil {
		return err
	}
	replicaCount := roundToInt(replicaCountFloat)

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	fn := func(object *unstructured.Unstructured) error {
		return unstructured.SetNestedField(object.Object, replicaCount, "spec", "replicas")
	}

	return e.store.Update(ctx, key, fn)
}
