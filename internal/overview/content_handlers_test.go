package overview

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	cachefake "github.com/heptio/developer-dash/internal/cache/fake"
	cacheutil "github.com/heptio/developer-dash/internal/cache/util"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/clock"
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
		init     func(*testing.T, *cachefake.MockCache)
		fields   map[string]string
		keys     []cacheutil.Key
		expected []*unstructured.Unstructured
		isErr    bool
	}{
		{
			name: "without name",
			init: func(t *testing.T, c *cachefake.MockCache) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				c.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(sampleObjects, nil)
			},
			fields:   map[string]string{},
			keys:     []cacheutil.Key{{APIVersion: "v1", Kind: "kind"}},
			expected: sortedSampleObjects,
		},
		{
			name: "name",
			init: func(t *testing.T, c *cachefake.MockCache) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind",
					Name:       "name"}

				c.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return([]*unstructured.Unstructured{}, nil)

			},
			fields: map[string]string{"name": "name"},
			keys:   []cacheutil.Key{{APIVersion: "v1", Kind: "kind"}},
		},
		{
			name: "cache retrieve error",
			init: func(t *testing.T, c *cachefake.MockCache) {
				key := cacheutil.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				c.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, errors.New("error"))
			},
			fields: map[string]string{},
			keys:   []cacheutil.Key{{APIVersion: "v1", Kind: "kind"}},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			c := cachefake.NewMockCache(controller)
			tc.init(t, c)

			namespace := "default"

			ctx := context.Background()
			got, err := loadObjects(ctx, c, namespace, tc.fields, tc.keys)
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
