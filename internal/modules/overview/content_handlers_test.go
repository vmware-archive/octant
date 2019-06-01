package overview

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/clock"

	"github.com/heptio/developer-dash/internal/describer"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/pkg/objectstoreutil"
)

func createObject(name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

func Test_loadObjects(t *testing.T) {
	sampleObjects := []*unstructured.Unstructured{
		createObject("omega"),
		createObject("alpha"),
	}
	sortedSampleObjects := []*unstructured.Unstructured{
		createObject("alpha"),
		createObject("omega"),
	}

	cases := []struct {
		name     string
		init     func(*testing.T, *storefake.MockObjectStore)
		fields   map[string]string
		keys     []objectstoreutil.Key
		expected []*unstructured.Unstructured
		isErr    bool
	}{
		{
			name: "without name",
			init: func(t *testing.T, o *storefake.MockObjectStore) {
				key := objectstoreutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(sampleObjects, nil)
			},
			fields:   map[string]string{},
			keys:     []objectstoreutil.Key{{APIVersion: "v1", Kind: "kind"}},
			expected: sortedSampleObjects,
		},
		{
			name: "name",
			init: func(t *testing.T, o *storefake.MockObjectStore) {
				key := objectstoreutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind",
					Name:       "name"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return([]*unstructured.Unstructured{}, nil)

			},
			fields: map[string]string{"name": "name"},
			keys:   []objectstoreutil.Key{{APIVersion: "v1", Kind: "kind"}},
		},
		{
			name: "cache retrieve error",
			init: func(t *testing.T, o *storefake.MockObjectStore) {
				key := objectstoreutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("error"))
			},
			fields: map[string]string{},
			keys:   []objectstoreutil.Key{{APIVersion: "v1", Kind: "kind"}},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storefake.NewMockObjectStore(controller)
			tc.init(t, o)

			namespace := "default"

			ctx := context.Background()
			got, err := describer.LoadObjects(ctx, o, namespace, tc.fields, tc.keys)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}

}

func Test_translateTimestamp(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	cases := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "zero",
			expected: "<unknown>",
		},
		{
			name:     "not zero",
			time:     time.Unix(1538828100, 0),
			expected: "30s",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := metav1.NewTime(tc.time)

			got := translateTimestamp(ts, c)
			assert.Equal(t, tc.expected, got)
		})
	}
}
