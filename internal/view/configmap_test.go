package view

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/kubernetes/pkg/apis/core"
)

func TestConfigMapDetails_InvalidObject(t *testing.T) {
	cm := NewConfigMapDetails()
	ctx := context.Background()

	object := &unstructured.Unstructured{}

	_, err := cm.Content(ctx, object, nil)
	require.Error(t, err)
}

func TestConfigMapDetails(t *testing.T) {
	cm := NewConfigMapDetails()

	ctx := context.Background()
	object := &core.ConfigMap{
		Data: map[string]string{
			"test": "data",
		},
	}

	contents, err := cm.Content(ctx, object, nil)
	require.NoError(t, err)

	require.Len(t, contents, 1)

	table, ok := contents[0].(*content.Table)
	require.True(t, ok)
	require.Len(t, table.Rows, 1)

	expectedColumns := []string{"Key", "Value"}
	assert.Equal(t, expectedColumns, table.ColumnNames())

	expectedRow := content.TableRow{
		"Key":   content.NewStringText("test"),
		"Value": content.NewStringText("data"),
	}
	assert.Equal(t, expectedRow, table.Rows[0])
}
