/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/pointer"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	actionFake "github.com/vmware-tanzu/octant/pkg/action/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDeploymentConfigurationEditor(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger := log.NopLogger()

	deployment := testutil.CreateDeployment("deployment")
	deployment.Namespace = "default"

	objectStore := fake.NewMockStore(controller)
	alerter := actionFake.NewMockAlerter(controller)

	key, err := store.KeyFromObject(deployment)
	require.NoError(t, err)

	updatedDeployment := deployment.DeepCopy()
	updatedDeployment.Spec.Replicas = pointer.Int32Ptr(5)

	objectStore.EXPECT().
		Update(gomock.Any(), key, gomock.Any()).
		DoAndReturn(func(ctx context.Context, key store.Key, fn func(object *unstructured.Unstructured) error) error {
			return nil
		})

	alerter.EXPECT().
		SendAlert(gomock.Any()).
		DoAndReturn(func(alert action.Alert) {
			assert.Equal(t, action.AlertTypeInfo, alert.Type)
			assert.Equal(t, `Updated Deployment "deployment"`, alert.Message)
			assert.NotNil(t, alert.Expiration)
		})

	configurationEditor := NewDeploymentConfigurationEditor(logger, objectStore)
	assert.Equal(t, "deployment/configuration", configurationEditor.ActionName())

	ctx := context.Background()

	payload := action.Payload{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"namespace":  "default",
		"name":       "deployment",
		"replicas":   "5",
	}

	require.NoError(t, configurationEditor.Handle(ctx, alerter, payload))

}
