package printer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
	"github.com/heptio/developer-dash/pkg/view/flexlayout"
)

func mockObjectsEvents(t *testing.T, appObjectStore *storefake.MockObjectStore, namespace string, events ...corev1.Event) {
	require.NotNil(t, appObjectStore)

	var objects []*unstructured.Unstructured

	for _, event := range events {
		objects = append(objects, testutil.ToUnstructured(t, &event))
	}

	key := objectstoreutil.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Event",
	}

	appObjectStore.EXPECT().
		List(gomock.Any(), key).
		Return(objects, nil).
		AnyTimes()
}

func stubMetadataForObject(t *testing.T, object runtime.Object, fl *flexlayout.FlexLayout) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	metadata, err := NewMetadata(object, tpo.link)
	require.NoError(t, err)
	require.NoError(t, metadata.AddToFlexLayout(fl))
}
