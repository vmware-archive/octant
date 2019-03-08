package printer

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/testutil"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func Test_CustomResourceListHandler(t *testing.T) {
	crd := loadCRDFromFile(t, "crd.yaml")
	resource := loadCRFromFile(t, "crd-resource.yaml")

	now := time.Now()
	resource.SetCreationTimestamp(metav1.Time{Time: now})

	labels := map[string]string{"foo": "bar"}
	resource.SetLabels(labels)

	list := []*unstructured.Unstructured{
		resource,
	}
	ctx := context.Background()
	got, err := CustomResourceListHandler(ctx, crd.Name, "default", crd, list)
	require.NoError(t, err)

	expected := component.NewTableWithRows(
		"crontabs.stable.example.com",
		component.NewTableCols("Name", "Labels", "Age"),
		[]component.TableRow{
			{
				"Name":   component.NewLink("", "my-crontab", "/content/overview/namespace/default/custom-resources/crontabs.stable.example.com/my-crontab"),
				"Age":    component.NewTimestamp(now),
				"Labels": component.NewLabels(labels),
			},
		})

	assert.Equal(t, expected, got)
}

func Test_CustomResourceListHandler_custom_columns(t *testing.T) {
	crd := loadCRDFromFile(t, "crd-additional-columns.yaml")
	resource := loadCRFromFile(t, "crd-resource.yaml")

	now := time.Now()
	resource.SetCreationTimestamp(metav1.Time{Time: now})

	labels := map[string]string{"foo": "bar"}
	resource.SetLabels(labels)

	list := []*unstructured.Unstructured{
		resource,
	}

	ctx := context.Background()
	got, err := CustomResourceListHandler(ctx, crd.Name, "default", crd, list)
	require.NoError(t, err)

	expected := component.NewTableWithRows(
		"crontabs.stable.example.com",
		component.NewTableCols("Name", "Labels", "Spec", "Replicas", "Errors", "Age"),
		[]component.TableRow{
			{
				"Name":     component.NewLink("", "my-crontab", "/content/overview/namespace/default/custom-resources/crontabs.stable.example.com/my-crontab"),
				"Age":      component.NewTimestamp(now),
				"Labels":   component.NewLabels(labels),
				"Replicas": component.NewText("1"),
				"Spec":     component.NewText("* * * * */5"),
				"Errors":   component.NewText("1"),
			},
		})

	assert.Equal(t, expected, got)
}

func Test_printResourceConfig(t *testing.T) {
	cases := []struct {
		name     string
		crd      string
		cr       string
		expected component.ViewComponent
		isErr    bool
	}{
		{
			name: "with additional columns",
			crd:  "crd-additional-columns.yaml",
			cr:   "crd-resource.yaml",
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Spec",
					Content: component.NewText("* * * * */5"),
				},
				{
					Header:  "Replicas",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name: "in general",
			crd:  "crd.yaml",
			cr:   "crd-resource.yaml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			crd := loadCRDFromFile(t, tc.crd)
			resource := loadCRFromFile(t, tc.cr)

			now := time.Now()
			resource.SetCreationTimestamp(metav1.Time{Time: now})

			labels := map[string]string{"foo": "bar"}
			resource.SetLabels(labels)

			got, err := printCustomResourceConfig(resource, crd)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_printResourceStatus(t *testing.T) {
	cases := []struct {
		name     string
		crd      string
		cr       string
		expected component.ViewComponent
		isErr    bool
	}{
		{
			name: "with additional columns",
			crd:  "crd-additional-columns.yaml",
			cr:   "crd-resource.yaml",
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Errors",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name: "in general",
			crd:  "crd.yaml",
			cr:   "crd-resource.yaml",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			crd := loadCRDFromFile(t, tc.crd)
			resource := loadCRFromFile(t, tc.cr)

			now := time.Now()
			resource.SetCreationTimestamp(metav1.Time{Time: now})

			labels := map[string]string{"foo": "bar"}
			resource.SetLabels(labels)

			got, err := printCustomResourceStatus(resource, crd)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func loadCRDFromFile(t *testing.T, filename string) *apiextv1beta1.CustomResourceDefinition {
	crd := testutil.CreateCRD("crd")
	testutil.LoadTypedObjectFromFile(t, filename, crd)

	return crd
}

func loadCRFromFile(t *testing.T, filename string) *unstructured.Unstructured {
	file, err := os.Open(filepath.Join("testdata", filename))
	require.NoError(t, err)

	decoder := yaml.NewYAMLOrJSONDecoder(file, 1024)
	var m map[string]interface{}
	require.NoError(t, decoder.Decode(&m))

	resource := &unstructured.Unstructured{
		Object: m,
	}

	return resource
}
