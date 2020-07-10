/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant_test

import (
	"context"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testClient "k8s.io/client-go/kubernetes/fake"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	actionFake "github.com/vmware-tanzu/octant/pkg/action/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
)

func Test_CronJobTrigger(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)
	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	fakeClientset := testClient.NewSimpleClientset()

	clusterClient.EXPECT().KubernetesClient().AnyTimes().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().BatchV1().AnyTimes().Return(fakeClientset.BatchV1())

	key := store.Key{
		Namespace:  "namespace",
		APIVersion: "v1beta1",
		Kind:       "CronJob",
		Name:       "cron",
	}

	cronjob := testutil.CreateCronJob("cron")

	objectStore.EXPECT().
		Get(ctx, gomock.Eq(key)).
		Return(testutil.ToUnstructured(t, cronjob), nil)

	alertType := action.AlertTypeInfo

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, alertType, alert.Type)
			assert.NotNil(t, alert.Expiration)
		})

	jobToCreate := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "manual-job",
			Namespace:   cronjob.Namespace,
			Annotations: cronjob.Annotations,
			Labels:      cronjob.Labels,
		},
		Spec: cronjob.Spec.JobTemplate.Spec,
	}

	_, err := fakeClientset.BatchV1().Jobs(cronjob.Namespace).Create(context.TODO(), jobToCreate, metav1.CreateOptions{})
	require.NoError(t, err)

	trigger := octant.NewCronJobTrigger(objectStore, clusterClient)
	assert.Equal(t, octant.ActionOverviewCronjob, trigger.ActionName())

	payload := action.CreatePayload(octant.ActionOverviewCronjob, map[string]interface{}{
		"namespace":  key.Namespace,
		"apiVersion": key.APIVersion,
		"kind":       key.Kind,
		"name":       key.Name,
	})

	require.NoError(t, trigger.Handle(ctx, alerter, payload))
}

func Test_CronJobHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)
	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
	clusterClient := clusterFake.NewMockClientInterface(controller)
	fakeClientset := testClient.NewSimpleClientset()

	clusterClient.EXPECT().KubernetesClient().AnyTimes().Return(kubernetesClient, nil)
	kubernetesClient.EXPECT().BatchV1().AnyTimes().Return(fakeClientset.BatchV1())

	cronjob := testutil.CreateCronJob("cron")
	key, err := store.KeyFromObject(cronjob)
	require.NoError(t, err)

	objectStore.EXPECT().
		Get(ctx, gomock.Eq(key)).
		Return(testutil.ToUnstructured(t, cronjob), nil).AnyTimes()

	objectStore.EXPECT().
		Update(gomock.Any(), key, gomock.Any()).Return(nil).AnyTimes()

	alertType := action.AlertTypeInfo

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, alertType, alert.Type)
			assert.NotNil(t, alert.Expiration)
		}).AnyTimes()

	suspend := octant.NewCronJobSuspend(objectStore, clusterClient)
	assert.Equal(t, octant.ActionOverviewSuspendCronjob, suspend.ActionName())

	suspendPayload := action.CreatePayload(octant.ActionOverviewSuspendCronjob, map[string]interface{}{
		"namespace":  key.Namespace,
		"apiVersion": key.APIVersion,
		"kind":       key.Kind,
		"name":       key.Name,
	})

	require.NoError(t, suspend.Handle(ctx, alerter, suspendPayload))

	resume := octant.NewCronJobResume(objectStore, clusterClient)
	assert.Equal(t, octant.ActionOverviewResumeCronjob, resume.ActionName())

	resumePayload := action.CreatePayload(octant.ActionOverviewSuspendCronjob, map[string]interface{}{
		"namespace":  key.Namespace,
		"apiVersion": key.APIVersion,
		"kind":       key.Kind,
		"name":       key.Name,
	})

	require.NoError(t, suspend.Handle(ctx, alerter, resumePayload))
}
