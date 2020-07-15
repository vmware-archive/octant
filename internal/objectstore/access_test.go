package objectstore

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	clusterfake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Test_accessCache(t *testing.T) {
	c := newAccessCache()

	key := AccessKey{
		Namespace: "test",
		Group:     "group",
		Resource:  "resource",
		Verb:      "list",
	}

	c.set(key, true)

	got, isFound := c.get(key)
	require.True(t, isFound)
	require.True(t, got)
}

func Test_ResourceAccess_HasAccess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	client := clusterfake.NewMockClientInterface(controller)
	namespaceClient := clusterfake.NewMockNamespaceInterface(controller)

	client.EXPECT().NamespaceClient().Return(namespaceClient, nil).MaxTimes(3)
	namespaces := []string{"test", ""}
	namespaceClient.EXPECT().Names().Return(namespaces, nil).MaxTimes(3)

	r := NewResourceAccess(client)

	scenarios := []struct {
		name        string
		resource    string
		key         store.Key
		setupAccess func()
		expectErr   bool
	}{
		{
			name:     "pods",
			resource: "pods",
			key: store.Key{
				APIVersion: "apps/v1",
				Kind:       "Pod",
			},
			setupAccess: func() {
				aKey := AccessKey{
					Namespace: "",
					Group:     "apps",
					Resource:  "pods",
					Verb:      "get",
				}
				r.Set(aKey, true)
			},
			expectErr: false,
		},
		{
			name:     "crds",
			resource: "customresourcedefinitions",
			key: store.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			setupAccess: func() {
				aKey := AccessKey{
					Namespace: "",
					Group:     "apiextensions.k8s.io",
					Resource:  "customresourcedefinitions",
					Verb:      "get",
				}
				r.Set(aKey, true)
			},
			expectErr: false,
		},
		{
			name:     "no access crds",
			resource: "customresourcedefinitions",
			key: store.Key{
				APIVersion: "apiextensions.k8s.io/v1beta1",
				Kind:       "CustomResourceDefinition",
			},
			setupAccess: func() {
				aKey := AccessKey{
					Namespace: "",
					Group:     "apiextensions.k8s.io",
					Resource:  "customresourcedefinitions",
					Verb:      "get",
				}
				r.Set(aKey, false)
			},
			expectErr: true,
		},
	}

	for i := range scenarios {
		ts := scenarios[i]
		t.Run(ts.name, func(t *testing.T) {
			ts.setupAccess()

			gvk := ts.key.GroupVersionKind()
			podGVR := schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: ts.resource,
			}
			client.EXPECT().Resource(gomock.Eq(gvk.GroupKind())).Return(podGVR, true, nil)

			if ts.expectErr {
				require.Error(t, r.HasAccess(ctx, ts.key, "get"))
			} else {
				require.NoError(t, r.HasAccess(ctx, ts.key, "get"))
			}
		})
	}
}
