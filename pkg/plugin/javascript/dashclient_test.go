package javascript

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
)

func TestDashboardClient_Get(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	objectStore := storeFake.NewMockStore(controller)

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	obj := CreateDashClientObject(ctx, objectStore, vm)
	assert.NotNil(t, obj)

	objectStore.EXPECT().Get(ctx, store.Key{APIVersion: "v1", Kind: "Pod"}).Return(&unstructured.Unstructured{}, nil)

	vm.Set("dashClient", obj)
	_, err := vm.RunString(`
dashClient.Get({apiVersion: 'v1', kind:'Pod'})
`)
	assert.NoError(t, err)
}

func TestDashboardClient_List(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	objectStore := storeFake.NewMockStore(controller)

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	obj := CreateDashClientObject(ctx, objectStore, vm)
	assert.NotNil(t, obj)

	objectStore.EXPECT().List(ctx, store.Key{Namespace: "test", APIVersion: "v1", Kind: "Pod"}).Return(&unstructured.UnstructuredList{}, false, nil)

	vm.Set("dashClient", obj)
	_, err := vm.RunString(`
dashClient.List({namespace:'test', apiVersion: 'v1', kind:'Pod'})
`)
	assert.NoError(t, err)
}

func TestDashboardClient_CreateUpdate(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	objectStore := storeFake.NewMockStore(controller)

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	obj := CreateDashClientObject(ctx, objectStore, vm)
	assert.NotNil(t, obj)

	objectStore.EXPECT().CreateOrUpdateFromYAML(ctx, "test", "create-yaml").Return([]string{}, nil)
	objectStore.EXPECT().CreateOrUpdateFromYAML(ctx, "test", "update-yaml").Return([]string{}, nil)

	vm.Set("dashClient", obj)
	_, err := vm.RunString(`
dashClient.Create('test', 'create-yaml')
dashClient.Update('test', 'update-yaml')
`)
	assert.NoError(t, err)
}

func TestDashboardClient_Delete(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()
	objectStore := storeFake.NewMockStore(controller)

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	obj := CreateDashClientObject(ctx, objectStore, vm)
	assert.NotNil(t, obj)

	objectStore.EXPECT().Delete(ctx, store.Key{Namespace: "test", APIVersion: "v1", Kind: "ReplicaSet", Name: "my-replica-set"}).Return(nil)

	vm.Set("dashClient", obj)
	_, err := vm.RunString(`
dashClient.Delete({namespace:'test', apiVersion: 'v1', kind:'ReplicaSet', name: 'my-replica-set'})
`)
	assert.NoError(t, err)
}
