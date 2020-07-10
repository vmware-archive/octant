/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// CronJobTrigger manually triggers a cronjob
type CronJobTrigger struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*CronJobTrigger)(nil)

// NewCronJobTrigger creates an instance of CronJobTrigger
func NewCronJobTrigger(objectStore store.Store, clusterClient cluster.ClientInterface) *CronJobTrigger {
	cronjob := &CronJobTrigger{
		store:         objectStore,
		clusterClient: clusterClient,
	}
	return cronjob
}

// ActionName returns the name of this action
func (c *CronJobTrigger) ActionName() string {
	return "action.octant.dev/cronJob"
}

// Handle executing cronjob
func (c *CronJobTrigger) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
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
		return errors.New("object store cannot get cronjob")
	}

	cronjob := &batchv1beta1.CronJob{}
	if err := kubernetes.FromUnstructured(object, cronjob); err != nil {
		return err
	}

	newJobName := createJobName(cronjob.Name)

	var message string
	alertType := action.AlertTypeInfo
	if err := c.Trigger(newJobName, cronjob); err != nil {
		message = fmt.Sprintf("Unable to create job %q: %s", key.Name, err)
		logger := log.From(ctx)
		logger.WithErr(err).Errorf("trigger cronjob")
	}
	message = fmt.Sprintf("Job %s created", newJobName)
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)
	return nil
}

// Trigger manually creates a new job
func (c *CronJobTrigger) Trigger(name string, cronJob *batchv1beta1.CronJob) error {
	if cronJob == nil {
		return errors.New("nil cronjob")
	}

	client, err := c.clusterClient.KubernetesClient()
	if err != nil {
		return err
	}

	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"

	labels := make(map[string]string)
	for k, v := range cronJob.Spec.JobTemplate.Labels {
		labels[k] = v
	}

	jobToCreate := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   cronJob.Namespace,
			Annotations: annotations,
			Labels:      labels,
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	_, err = client.BatchV1().Jobs(cronJob.Namespace).Create(context.TODO(), jobToCreate, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// createJobName creates a job name
func createJobName(cronJobName string) string {
	// From https://github.com/kubernetes/dashboard/blob/v2.0.0-rc5/src/app/backend/resource/cronjob/jobs.go#L81
	var newJobName string
	if len(cronJobName) < 42 {
		newJobName = cronJobName + "-manual-" + rand.String(3)
	} else {
		newJobName = cronJobName[0:41] + "-manual-" + rand.String(3)
	}

	return newJobName
}

// CronJobSuspend pauses a cronjob
type CronJobSuspend struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*CronJobSuspend)(nil)

func NewCronJobSuspend(objectStore store.Store, clusterClient cluster.ClientInterface) *CronJobSuspend {
	cronjob := &CronJobSuspend{
		store:         objectStore,
		clusterClient: clusterClient,
	}
	return cronjob
}

// Handle suspending cronjob
func (c *CronJobSuspend) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", c.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	expiration := time.Now().Add(10 * time.Second)

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	object, err := c.store.Get(ctx, key)
	if err != nil {
		return err
	}

	if object == nil {
		return errors.New("object store cannot get cronjob")
	}

	cronjob := &batchv1beta1.CronJob{}
	if err := kubernetes.FromUnstructured(object, cronjob); err != nil {
		return err
	}

	if cronjob.Spec.Suspend != nil {
		*cronjob.Spec.Suspend = true
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cronjob)
	if err != nil {
		return nil
	}

	unstructuredCronJob := &unstructured.Unstructured{Object: m}

	err = c.store.Update(ctx, key, func(u *unstructured.Unstructured) error {
		if unstructuredCronJob.GetAPIVersion() != u.GetAPIVersion() {
			return fmt.Errorf("object API version cannot be updated")
		}
		if unstructuredCronJob.GetKind() != u.GetKind() {
			return fmt.Errorf("object kind cannot be updated")
		}
		if unstructuredCronJob.GetName() != u.GetName() {
			return fmt.Errorf("object name cannot be updated")
		}

		delete(unstructuredCronJob.Object, "status")

		for k := range unstructuredCronJob.Object {
			u.Object[k] = unstructuredCronJob.Object[k]
		}
		return nil
	})

	if err != nil {
		sendAlert(
			alerter,
			action.AlertTypeError,
			fmt.Sprintf("update: %s", err.Error()),
			&expiration,
		)
	}

	successMessage := fmt.Sprintf("Suspending %s (%s) %s in %s",
		unstructuredCronJob.GetKind(),
		unstructuredCronJob.GetAPIVersion(),
		unstructuredCronJob.GetName(),
		unstructuredCronJob.GetNamespace())
	sendAlert(alerter, action.AlertTypeInfo, successMessage, &expiration)
	return nil
}

// ActionName returns the action name
func (c *CronJobSuspend) ActionName() string {
	return ActionOverviewSuspendCronjob
}

// CronJobResume resumes a cronjob
type CronJobResume struct {
	store         store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*CronJobResume)(nil)

func NewCronJobResume(objectStore store.Store, clusterClient cluster.ClientInterface) *CronJobResume {
	cronjob := &CronJobResume{
		store:         objectStore,
		clusterClient: clusterClient,
	}
	return cronjob
}

// Handle resuming cronjob
func (c *CronJobResume) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	logger := log.From(ctx).With("actionName", c.ActionName())
	logger.With("payload", payload).Infof("received action payload")

	expiration := time.Now().Add(10 * time.Second)

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	object, err := c.store.Get(ctx, key)
	if err != nil {
		return err
	}

	if object == nil {
		return errors.New("object store cannot get cronjob")
	}

	cronjob := &batchv1beta1.CronJob{}
	if err := kubernetes.FromUnstructured(object, cronjob); err != nil {
		return err
	}

	*cronjob.Spec.Suspend = false

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cronjob)
	if err != nil {
		return nil
	}

	unstructuredCronJob := &unstructured.Unstructured{Object: m}

	err = c.store.Update(ctx, key, func(u *unstructured.Unstructured) error {
		if unstructuredCronJob.GetAPIVersion() != u.GetAPIVersion() {
			return fmt.Errorf("object API version cannot be updated")
		}
		if unstructuredCronJob.GetKind() != u.GetKind() {
			return fmt.Errorf("object kind cannot be updated")
		}
		if unstructuredCronJob.GetName() != u.GetName() {
			return fmt.Errorf("object name cannot be updated")
		}

		delete(unstructuredCronJob.Object, "status")

		for k := range unstructuredCronJob.Object {
			u.Object[k] = unstructuredCronJob.Object[k]
		}
		return nil
	})

	if err != nil {
		sendAlert(
			alerter,
			action.AlertTypeError,
			fmt.Sprintf("update: %s", err.Error()),
			&expiration,
		)
	}

	successMessage := fmt.Sprintf("Resuming %s (%s) %s in %s",
		unstructuredCronJob.GetKind(),
		unstructuredCronJob.GetAPIVersion(),
		unstructuredCronJob.GetName(),
		unstructuredCronJob.GetNamespace())
	sendAlert(alerter, action.AlertTypeInfo, successMessage, &expiration)
	return nil
}

// ActionName returns the action name
func (c *CronJobResume) ActionName() string {
	return ActionOverviewResumeCronjob
}
