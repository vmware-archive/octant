package overview

import (
	"context"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestIngressSummary_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressSummary("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestIngressDetails_InvalidObject(t *testing.T) {
	assertViewInvalidObject(t, NewIngressDetails("prefix", "ns", clock.NewFakeClock(time.Now())))
}

func TestIngressDetails(t *testing.T) {
	v := NewIngressDetails("prefix", "ns", clock.NewFakeClock(time.Now()))

	c := cache.NewMemoryCache()

	ingress := loadFromFile(t, "ingress-1.yaml")

	ctx := context.Background()

	got, err := v.Content(ctx, ingress, c)
	require.NoError(t, err)

	tlsTable := content.NewTable("TLS", "TLS is not configured for this Ingress")
	tlsTable.Columns = view.TableCols("Secret", "Hosts")

	rulesTable := content.NewTable("Rules", "Rules are not configured for this Ingress")
	rulesTable.Columns = view.TableCols("Host", "Path", "Backend")
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
