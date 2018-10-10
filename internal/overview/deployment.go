package overview

import (
	"fmt"
	"net/url"
	"path"

	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

// DeploymentsDescriber creates content for a list of deployments.
type DeploymentsDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

// NewDeploymentsDescriber creates an instance of DeploymentsDescriber.
func NewDeploymentsDescriber() *DeploymentsDescriber {
	return &DeploymentsDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
		},
	}
}

// Describe creates content.
func (d *DeploymentsDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	var contents []Content

	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	if len(objects) < 1 {
		return contents, nil
	}

	t := newDeploymentTable("Deployments")
	for _, object := range objects {
		cur := &appsv1.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, cur)
		if err != nil {
			return nil, err
		}

		t.Rows = append(t.Rows, printDeployment(cur, prefix, namespace, d.clock()))
	}

	contents = append(contents, t)

	return contents, nil
}

// DeploymentDescriber creates content for a single deployment.
type DeploymentDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

// NewDeploymentDescriber creates an instance of DeploymentDescriber.
func NewDeploymentDescriber() *DeploymentDescriber {
	return &DeploymentDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
			},
		},
	}
}

// Describe creates content.
func (d *DeploymentDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	var contents []Content

	t := newDeploymentTable("Deployment")

	if len(objects) != 1 {
		return nil, errors.Errorf("expected 1 deployment")
	}

	object := objects[0]
	obj := &appsv1.Deployment{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, obj)
	if err != nil {
		return nil, err
	}

	t.Rows = append(t.Rows, printDeployment(obj, prefix, namespace, d.clock()))

	contents = append(contents, t)

	eventsTable, err := eventsForObject(object, cache, prefix, namespace, d.clock())
	if err != nil {
		return nil, err
	}

	contents = append(contents, eventsTable)

	return contents, nil
}

func newDeploymentTable(name string) table {
	t := newTable(name)

	t.Columns = []tableColumn{
		{Name: "Name", Accessor: "name"},
		{Name: "Labels", Accessor: "labels"},
		{Name: "Pods", Accessor: "pods"},
		{Name: "Age", Accessor: "age"},
		{Name: "Images", Accessor: "images"},
	}

	return t
}

func printDeployment(d *appsv1.Deployment, prefix, namespace string, c clock.Clock) tableRow {
	var images []string
	for _, container := range d.Spec.Template.Spec.Containers {
		images = append(images, container.Image)
	}

	pods := fmt.Sprintf("%d/%d",
		d.Status.Replicas,
		*d.Spec.Replicas,
	)

	values := url.Values{}
	values.Set("namespace", namespace)

	deploymentPath := fmt.Sprintf("%s?%s",
		path.Join(prefix, "/workloads/deployments", d.GetName()),
		values.Encode())

	return tableRow{
		"name":   newLinkText(d.GetName(), deploymentPath),
		"labels": newLabelsText(d.GetLabels()),
		"pods":   newStringText(pods),
		"age":    newStringText(translateTimestamp(d.CreationTimestamp, c)),
		"images": newListText(images),
	}
}
