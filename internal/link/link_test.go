/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package link

import (
	"net/url"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/link/fake"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func TestLink_ForObject(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	controller := gomock.NewController(t)
	defer controller.Finish()

	config := fake.NewMockConfig(controller)

	config.EXPECT().
		ObjectPath(
			gomock.Eq(deployment.Namespace),
			gomock.Eq(deployment.APIVersion),
			gomock.Eq(deployment.Kind),
			gomock.Eq(deployment.Name)).
		Return("/path", nil)

	l, err := NewFromDashConfig(config)
	require.NoError(t, err)

	got, err := l.ForObject(deployment, "my object")
	require.NoError(t, err)

	expectedRef := path.Join("/path")
	assert.Equal(t, expectedRef, got.Ref())
	assert.Equal(t, "my object", got.Text())
}

func TestLink_ForObjectWithQuery(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	controller := gomock.NewController(t)
	defer controller.Finish()

	config := fake.NewMockConfig(controller)

	config.EXPECT().
		ObjectPath(
			gomock.Eq(deployment.Namespace),
			gomock.Eq(deployment.APIVersion),
			gomock.Eq(deployment.Kind),
			gomock.Eq(deployment.Name)).
		Return("/path", nil)

	l, err := NewFromDashConfig(config)
	require.NoError(t, err)

	query := url.Values{}
	query.Set("foo", "bar")
	got, err := l.ForObjectWithQuery(deployment, "my object", query)
	require.NoError(t, err)

	p := path.Join("/path")
	u := url.URL{Path: p, RawQuery: query.Encode()}
	assert.Equal(t, u.String(), got.Ref())
	assert.Equal(t, "my object", got.Text())
}

func TestLink_ForGVK(t *testing.T) {
	namespace := "default"
	apiVersion := "v1"
	kind := "Pod"
	name := "pod"
	text := "pod"

	controller := gomock.NewController(t)
	defer controller.Finish()

	config := fake.NewMockConfig(controller)

	config.EXPECT().
		ObjectPath(
			gomock.Eq("default"),
			gomock.Eq("v1"),
			gomock.Eq("Pod"),
			gomock.Eq("pod")).
		Return("/path", nil)

	l, err := NewFromDashConfig(config)
	require.NoError(t, err)

	got, err := l.ForGVK(namespace, apiVersion, kind, name, text)
	require.NoError(t, err)

	expectedRef := path.Join("/path")
	assert.Equal(t, expectedRef, got.Ref())
	assert.Equal(t, "pod", got.Text())
}

func TestLink_ForOwner(t *testing.T) {
	deployment := testutil.CreateDeployment("deployment")

	ownerReference := &metav1.OwnerReference{
		APIVersion: "apiVersion",
		Kind:       "kind",
		Name:       "name",
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	config := fake.NewMockConfig(controller)

	config.EXPECT().
		ObjectPath(
			gomock.Eq("namespace"),
			gomock.Eq("apiVersion"),
			gomock.Eq("kind"),
			gomock.Eq("name")).
		Return("/path", nil)

	l, err := NewFromDashConfig(config)
	require.NoError(t, err)

	got, err := l.ForOwner(deployment, ownerReference)
	require.NoError(t, err)

	expectedRef := path.Join("/path")
	assert.Equal(t, expectedRef, got.Ref())
	assert.Equal(t, "name", got.Text())
}
