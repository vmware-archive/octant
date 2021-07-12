package overview

import (
	"context"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/cache"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/describer"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/config"
	"github.com/vmware-tanzu/octant/pkg/icon"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/store"
	objectStoreFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func Test_Overview_CRD_Navigation_Cache(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	crdKey := store.Key{
		APIVersion: "apiextensions.k8s.io/v1",
		Kind:       "CustomResourceDefinition",
	}

	client := clusterFake.NewMockClientInterface(controller)
	objectStore := objectStoreFake.NewMockStore(controller)
	objectStore.EXPECT().
		Watch(ctx, crdKey, gomock.Any()).
		DoAndReturn(func(_ context.Context, key store.Key, c *cache.ResourceEventHandlerFuncs) error {
			assert.NotNil(t, c.AddFunc)
			assert.NotNil(t, c.DeleteFunc)
			return nil
		})

	dashConfig := configFake.NewMockDash(controller)
	dashConfig.EXPECT().ObjectStore().Return(objectStore).AnyTimes()

	crdWatcher, err := describer.NewDefaultCRDWatcher(ctx, client, objectStore, nil)
	require.NoError(t, err)

	watchConfig := &config.CRDWatchConfig{
		Add:          func(ctx context.Context, object *unstructured.Unstructured) {},
		Delete:       func(ctx context.Context, object *unstructured.Unstructured) {},
		IsNamespaced: false,
	}

	require.NoError(t, crdWatcher.AddConfig(watchConfig))
	require.NoError(t, crdWatcher.Watch(ctx))

	dashConfig.EXPECT().Validate().Return(nil).AnyTimes()
	dashConfig.EXPECT().Logger().Return(log.NopLogger()).AnyTimes()
	dashConfig.EXPECT().CRDWatcher().Return(crdWatcher).AnyTimes()

	overviewOptions := Options{
		Namespace:  "test-namespace",
		DashConfig: dashConfig,
	}

	overviewModule, err := New(ctx, overviewOptions)
	require.NoError(t, err)
	assert.NotNil(t, overviewModule)
	assert.Nil(t, overviewModule.navigationCrdCache)

	crd := testutil.CreateCRDWithKind("namespace-scoped", "NamespaceScoped", false)
	crds := testutil.ToUnstructuredList(t, crd)
	objectStore.EXPECT().
		List(gomock.Any(), crdKey).
		Return(crds, false, nil).
		AnyTimes()

	namespaceName := "test"
	crNamespaceKey := store.Key{
		Namespace:  namespaceName,
		APIVersion: "testing/v1",
		Kind:       "NamespaceScoped",
	}

	cr := testutil.CreateCR("testing", "v1", "NamespaceScoped", "namespace-scoped")
	crs := testutil.ToUnstructuredList(t, cr)

	objectStore.EXPECT().
		List(gomock.Any(), crNamespaceKey).
		Return(crs, false, nil).
		AnyTimes()

	nav, err := navigation.New(cr.GetName(), path.Join("/prefix", cr.GetName()), navigation.SetNavigationIcon(icon.CustomResourceDefinition))
	require.NoError(t, err)

	namespaceGot, _, err := overviewModule.CRDEntries(ctx, "/prefix", namespaceName, objectStore, false)
	require.NoError(t, err)

	namespaceExpected := []navigation.Navigation{*nav}
	assert.Equal(t, namespaceExpected, namespaceGot)
	assert.NotNil(t, overviewModule.navigationCrdCache)

	expectedCache := map[string][]navigation.Navigation{namespaceName: namespaceExpected}
	assert.Equal(t, expectedCache, overviewModule.navigationCrdCache)
}
