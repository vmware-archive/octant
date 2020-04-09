/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// Cordon cordons a node
type Cordon struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*Cordon)(nil)

// NewCordon creates an instance of Cordon
func NewCordon(objectStore store.Store, clusterClient cluster.ClientInterface) *Cordon {
	cordon := &Cordon{
		store:         objectStore,
		clusterClient: clusterClient,
	}

	return cordon
}

// ActionName returns the name of this action
func (c *Cordon) ActionName() string {
	return "action.octant.dev/cordon"
}

// Handle executing cordon
func (c *Cordon) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", c.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	object, err := c.store.Get(ctx, key)
	if err != nil {
		return err
	}

	if object == nil {
		return errors.New("object store cannot get node")
	}

	var node *corev1.Node
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &node); err != nil {
		return err
	}

	message := fmt.Sprintf("Node %q marked as unschedulable", key.Name)
	alertType := action.AlertTypeInfo
	if err := c.Cordon(node); err != nil {
		message = fmt.Sprintf("Unable to cordon node %q: %s", key.Name, err)
		alertType = action.AlertTypeWarning
		logger := log.From(ctx)
		logger.WithErr(err).Errorf("cordon node")
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)
	return nil
}

// Cordon marks a node as unschedulable
func (c *Cordon) Cordon(node *corev1.Node) error {
	if node == nil {
		return errors.New("nil node")
	}

	client, err := c.clusterClient.KubernetesClient()
	if err != nil {
		return err
	}

	currentNode, err := client.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to find node %q", node.Name)
	}

	originalNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	if currentNode.Spec.Unschedulable {
		message := fmt.Sprintf("node %q already marked", node.Name)
		return errors.New(message)
	}
	currentNode.Spec.Unschedulable = true

	modifiedNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	patchBytes, patchErr := strategicpatch.CreateTwoWayMergePatch(originalNode, modifiedNode, node)
	if patchErr != nil {
		_, err = client.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchBytes)
	} else {
		_, err = client.CoreV1().Nodes().Update(currentNode)
		return errors.Wrapf(err, "failed to cordon %q", node.Name)
	}

	return err
}

// Uncordon uncordons a node
type Uncordon struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*Uncordon)(nil)

// NewUncordon creates an instances of uncordon
func NewUncordon(objectStore store.Store, clusterClient cluster.ClientInterface) *Uncordon {
	uncordon := &Uncordon{
		store:         objectStore,
		clusterClient: clusterClient,
	}

	return uncordon
}

// ActionName returns the name of this action
func (u *Uncordon) ActionName() string {
	return "action.octant.dev/uncordon"
}

// Handle executing uncordon
func (u *Uncordon) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", u.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	object, err := u.store.Get(ctx, key)
	if err != nil {
		return err
	}

	if object == nil {
		return errors.New("object store cannot get node")
	}

	var node *corev1.Node
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &node); err != nil {
		return err
	}

	message := fmt.Sprintf("Node %q marked as schedulable", key.Name)
	alertType := action.AlertTypeInfo
	if err := u.Uncordon(node); err != nil {
		message = fmt.Sprintf("Unable to uncordon node %q: %s", key.Name, err)
		alertType = action.AlertTypeWarning
		logger := log.From(ctx)
		logger.WithErr(err).Errorf("uncordon node")
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)
	return nil
}

// Uncordon marks a node as schedulable
func (u *Uncordon) Uncordon(node *corev1.Node) error {
	if node == nil {
		return errors.New("nil node")
	}

	client, err := u.clusterClient.KubernetesClient()
	if err != nil {
		return err
	}

	currentNode, err := client.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "unable to find node %q", node.Name)
	}

	originalNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	if !currentNode.Spec.Unschedulable {
		message := fmt.Sprintf("node %q already unmarked", node.Name)
		return errors.New(message)
	}
	currentNode.Spec.Unschedulable = false

	modifiedNode, err := json.Marshal(currentNode)
	if err != nil {
		return err
	}

	patchBytes, patchErr := strategicpatch.CreateTwoWayMergePatch(originalNode, modifiedNode, node)
	if patchErr != nil {
		_, err = client.CoreV1().Nodes().Patch(node.Name, types.StrategicMergePatchType, patchBytes)
	} else {
		_, err = client.CoreV1().Nodes().Update(currentNode)
		return errors.Wrapf(err, "failed to uncordon %q", node.Name)
	}

	return err
}
