package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/cluster/fake"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestClusterOverview(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy"),
	}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	o, err := NewClusterOverview(clusterClient, "default", log.NopLogger())
	require.NoError(t, err)
	if o == nil {
		return
	}

	assert.Equal(t, "overview", o.Name())
	assert.Equal(t, "/overview", o.ContentPath())
}

func TestClusterOverview_SetNamespace(t *testing.T) {
	scheme := runtime.NewScheme()
	objects := []runtime.Object{
		newUnstructured("apps/v1", "Deployment", "default", "deploy"),
	}

	clusterClient, err := fake.NewClient(scheme, resources, objects)
	require.NoError(t, err)

	o, err := NewClusterOverview(clusterClient, "default", log.NopLogger())
	require.NoError(t, err)
	if o == nil {
		return
	}
	defer o.Stop()

	err = o.SetNamespace("ns2")
	require.NoError(t, err)
}

// generatorFunc allows a bare Generate function to implement Generator
type generatorFunc func(ctx context.Context, path, prefix, namespace string) (component.ContentResponse, error)

func (g generatorFunc) Generate(ctx context.Context, path, prefix, namespace string) (component.ContentResponse, error) {
	return g(ctx, path, prefix, namespace)
}
