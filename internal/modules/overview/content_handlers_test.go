/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

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

	"github.com/vmware-tanzu/octant/internal/describer"
	internalErr "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/store"
	storeFake "github.com/vmware-tanzu/octant/pkg/store/fake"
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
	sampleObjects := testutil.ToUnstructuredList(t,
		createObject("omega"),
		createObject("alpha"))
	sortedSampleObjects := testutil.ToUnstructuredList(t,
		createObject("alpha"),
		createObject("omega"))

	cases := []struct {
		name     string
		init     func(*testing.T, *storeFake.MockStore)
		fields   map[string]string
		keys     []store.Key
		expected *unstructured.UnstructuredList
		isErr    bool
	}{
		{
			name: "without name",
			init: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(sampleObjects, false, nil)
			},
			fields:   map[string]string{},
			keys:     []store.Key{{APIVersion: "v1", Kind: "kind"}},
			expected: sortedSampleObjects,
		},
		{
			name: "name",
			init: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind",
					Name:       "name"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(&unstructured.UnstructuredList{}, false, nil)

			},
			fields:   map[string]string{"name": "name"},
			keys:     []store.Key{{APIVersion: "v1", Kind: "kind"}},
			expected: &unstructured.UnstructuredList{},
		},
		{
			name: "cache retrieve error",
			init: func(t *testing.T, o *storeFake.MockStore) {
				key := store.Key{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "kind"}

				o.EXPECT().
					List(gomock.Any(), gomock.Eq(key)).
					Return(nil, false, errors.New("error"))
			},
			fields: map[string]string{},
			keys:   []store.Key{{APIVersion: "v1", Kind: "kind"}},
			isErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			o := storeFake.NewMockStore(controller)
			tc.init(t, o)

			errorStore, err := internalErr.NewErrorStore()
			require.NoError(t, err)

			namespace := "default"

			ctx := context.Background()
			got, err := describer.LoadObjects(ctx, o, errorStore, namespace, tc.fields, tc.keys)
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
