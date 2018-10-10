package overview

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestDeploymentsDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "deployment.yaml")

	d := NewDeploymentsDescriber()
	got, err := d.Describe("/prefix", namespace, cache, nil)
	require.NoError(t, err)

	require.Len(t, got, 1)
	tbl, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, tbl.Title, "Deployments")
	assert.Len(t, tbl.Rows, 1)
}

func TestDeploymentDescriber(t *testing.T) {
	namespace := "default"

	cache := NewMemoryCache()
	loadUnstructured(t, cache, namespace, "deployment.yaml")

	fields := map[string]string{
		"name": "nginx-deployment",
	}

	d := NewDeploymentDescriber()
	got, err := d.Describe("/prefix", namespace, cache, fields)
	require.NoError(t, err)

	require.Len(t, got, 2)
	cjTable, ok := got[0].(table)
	require.True(t, ok)

	assert.Equal(t, cjTable.Title, "Deployment")
	assert.Len(t, cjTable.Rows, 1)

	eventsTable, ok := got[1].(table)
	require.True(t, ok)

	assert.Equal(t, eventsTable.Title, "Events")
	assert.Len(t, eventsTable.Rows, 0)
}

func Test_printDeployment(t *testing.T) {
	ti := time.Unix(1538828130, 0)
	c := clock.NewFakeClock(ti)

	d, ok := loadType(t, "deployment.yaml").(*appsv1.Deployment)
	require.True(t, ok)

	got := printDeployment(d, "/api", "default", c)

	expected := tableRow{
		"name":   newLinkText("nginx-deployment", "/api/workloads/deployments/nginx-deployment?namespace=default"),
		"labels": newLabelsText(map[string]string{"app": "nginx"}),
		"pods":   newStringText("3/3"),
		"age":    newStringText("10d"),
		"images": newListText([]string{"nginx:1.7.9"}),
	}

	assert.Equal(t, expected, got)
}
