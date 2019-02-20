package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ToUnstructured converts an object to an unstructured.
func ToUnstructured(t *testing.T, object runtime.Object) *unstructured.Unstructured {
	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	require.NoError(t, err)

	return &unstructured.Unstructured{Object: m}
}
