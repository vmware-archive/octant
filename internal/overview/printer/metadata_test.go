package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/overview/objectvisitor"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func Test_Metadata(t *testing.T) {
	fl := flexlayout.New()

	deployment := createDeployment("deployment")
	metadata, err := printer.NewMetadata(deployment)
	require.NoError(t, err)

	require.NoError(t, metadata.AddToFlexLayout(fl))

	got := fl.ToComponent("Summary")

	expected := &component.FlexLayout{
		Metadata: component.Metadata{
			Title: []component.TitleViewComponent{component.NewText("Summary")},
			Type:  "flexlayout",
		},
		Config: component.FlexLayoutConfig{
			Sections: []component.FlexLayoutSection{
				{
					{
						Width: 16,
						View: component.NewSummary("Metadata", component.SummarySections{
							{
								Header:  "Age",
								Content: component.NewTimestamp(deployment.CreationTimestamp.Time),
							},
						}...),
					},
				},
			},
		},
	}

	assert.Equal(t, expected, got)
}

func createDeployment(name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta:   genTypeMeta(objectvisitor.DeploymentGVK),
		ObjectMeta: genObjectMeta(name),
	}
}

func genTypeMeta(gvk schema.GroupVersionKind) metav1.TypeMeta {
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return metav1.TypeMeta{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}

func genObjectMeta(name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:              name,
		Namespace:         "namespace",
		UID:               types.UID(name),
		CreationTimestamp: metav1.Time{Time: time.Now()},
	}
}
