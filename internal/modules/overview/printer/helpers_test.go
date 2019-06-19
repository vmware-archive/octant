/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	configFake "github.com/vmware/octant/internal/config/fake"
	linkFake "github.com/vmware/octant/internal/link/fake"
	objectStoreFake "github.com/vmware/octant/pkg/store/fake"
	pluginFake "github.com/vmware/octant/pkg/plugin/fake"
	"github.com/vmware/octant/pkg/view/component"
)

const (
	rbacAPIVersion = "rbac.authorization.k8s.io/v1"
)

type testPrinterOptions struct {
	dashConfig *configFake.MockDash
	link       *linkFake.MockInterface

	objectStore   *objectStoreFake.MockStore
	pluginManager *pluginFake.MockManagerInterface
}

func newTestPrinterOptions(controller *gomock.Controller) *testPrinterOptions {
	objectStore := objectStoreFake.NewMockStore(controller)

	pluginManager := pluginFake.NewMockManagerInterface(controller)

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()

	tpo := &testPrinterOptions{
		dashConfig:    dashConfig,
		link:          linkFake.NewMockInterface(controller),
		objectStore:   objectStore,
		pluginManager: pluginManager,
	}

	tpo.dashConfig.EXPECT().Validate().Return(nil).AnyTimes()

	return tpo
}

func (o *testPrinterOptions) ToOptions() Options {
	return Options{
		DashConfig: o.dashConfig,
		Link:       o.link,
	}
}

func (o *testPrinterOptions) PathForObject(object runtime.Object, text, ref string) {
	l := component.NewLink("", text, ref)
	o.link.EXPECT().ForObject(gomock.Eq(object), text).Return(l, nil).AnyTimes()
}

func (o *testPrinterOptions) PathForGVK(namespace, apiVersion, kind, name, text, ref string) {
	l := component.NewLink("", text, ref)
	o.link.EXPECT().ForGVK(namespace, apiVersion, kind, name, text).Return(l, nil).AnyTimes()
}

func (o *testPrinterOptions) PathForOwner(parent runtime.Object, ownerReference *metav1.OwnerReference, ref string) {
	l := component.NewLink("", ownerReference.Name, ref)
	o.link.EXPECT().ForOwner(parent, ownerReference).Return(l, nil)
}

func assertComponentEqual(t *testing.T, expected, got component.Component) {
	a, err := json.Marshal(expected)
	require.NoError(t, err)

	b, err := json.Marshal(got)
	require.NoError(t, err)

	assert.JSONEq(t, string(a), string(b))
}
