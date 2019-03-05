package testutil

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

// LoadObjectFromFile loads a file from the `testdata` directory. It will
// assign it a `default` namespace if one is not set.
func LoadObjectFromFile(t *testing.T, objectFile string) runtime.Object {
	data, err := ioutil.ReadFile(filepath.Join("testdata", objectFile))
	require.NoError(t, err)

	object, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	require.NoError(t, err)

	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	require.NoError(t, err)
	if namespace == "" {
		require.NoError(t, accessor.SetNamespace(object, "default"))
	}

	return object
}

// LoadTypedObjectFromFile loads a file from the `testdata` directory. It will
// assign it a `default` namespace if one is not set.
func LoadTypedObjectFromFile(t *testing.T, objectFile string, into runtime.Object) {
	data, err := ioutil.ReadFile(filepath.Join("testdata", objectFile))
	require.NoError(t, err)

	gvk := into.GetObjectKind().GroupVersionKind()
	object, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, &gvk, into)
	require.NoError(t, err)

	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	require.NoError(t, err)
	if namespace == "" {
		require.NoError(t, accessor.SetNamespace(object, "default"))
	}
}
