/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	linkFake "github.com/vmware-tanzu/octant/internal/link/fake"
	portForwardFake "github.com/vmware-tanzu/octant/internal/portforward/fake"
	pluginFake "github.com/vmware-tanzu/octant/pkg/plugin/fake"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

const (
	rbacAPIVersion = "rbac.authorization.k8s.io/v1"
)

type testPrinterOptions struct {
	dashConfig *configFake.MockDash
	link       *linkFake.MockInterface

	objectStore   *objectStoreFake.MockStore
	pluginManager *pluginFake.MockManagerInterface
	portForwarder *portForwardFake.MockPortForwarder
}

func newTestPrinterOptions(controller *gomock.Controller) *testPrinterOptions {
	objectStore := objectStoreFake.NewMockStore(controller)

	pluginManager := pluginFake.NewMockManagerInterface(controller)

	portForwarder := portForwardFake.NewMockPortForwarder(controller)

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()
	dashConfig.EXPECT().PluginManager().Return(pluginManager).AnyTimes()
	dashConfig.EXPECT().PortForwarder().Return(portForwarder).AnyTimes()

	tpo := &testPrinterOptions{
		dashConfig:    dashConfig,
		link:          linkFake.NewMockInterface(controller),
		objectStore:   objectStore,
		pluginManager: pluginManager,
		portForwarder: portForwarder,
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

func gridActionsFactory(actions []component.GridAction) *component.GridActions {
	ga := component.NewGridActions()
	for _, gridAction := range actions {
		ga.AddGridAction(gridAction)
	}

	return ga
}

func buildObjectDeleteAction(t *testing.T, object runtime.Object) component.GridAction {
	action, err := objectDeleteAction(object)
	require.NoError(t, err)

	return action
}

// genObjectStatus generates object status for a link. It can be used
// when testing list handlers. This will be needed until there is a
// way to test that the list handlers are working without external
// dependencies.
func genObjectStatus(status component.TextStatus, messages []string) component.LinkOption {
	list := component.NewList(nil, nil)
	for _, msg := range messages {
		list.Add(component.NewText(msg))
	}

	return func(l *component.Link) {
		l.SetStatus(status, list)
	}
}
