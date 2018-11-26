package overview

import (
	"context"
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngressSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressSummary())
}

func TestIngressDetails_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressDetails())
}

func TestIngressDetails(t *testing.T) {
	v := NewIngressDetails()

	cache := NewMemoryCache()

	ingress := loadFromFile(t, "ingress-1.yaml")
	ingress = convertToInternal(t, ingress)

	ctx := context.Background()

	got, err := v.Content(ctx, ingress, cache)
	require.NoError(t, err)

	tlsTable := content.NewTable("TLS", "TLS is not configured for this Ingress")
	tlsTable.Columns = tableCols("Secret", "Hosts")

	rulesTable := content.NewTable("Rules", "Rules are not configured for this Ingress")
	rulesTable.Columns = tableCols("Host", "Path", "Backend")
	rulesTable.AddRow(content.TableRow{
		"Backend": content.NewLinkText("test:80", "/content/overview/discovery-and-load-balancing/services/test"),
		"Host":    content.NewStringText(""),
		"Path":    content.NewStringText("/testpath"),
	})

	expected := []content.Content{
		&tlsTable,
		&rulesTable,
	}

	assert.Equal(t, expected, got)
}
