package objectvisitor_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func sortObjectsByName(t *testing.T, objects []runtime.Object) {
	accessor := meta.NewAccessor()
	sort.Slice(objects, func(i, j int) bool {
		a, err := accessor.Name(objects[i])
		require.NoError(t, err)

		b, err := accessor.Name(objects[j])
		require.NoError(t, err)

		return a < b
	})

}
