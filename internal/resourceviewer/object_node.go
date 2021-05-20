package resourceviewer

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type objectNode struct {
	link          link.Interface
	pluginPrinter plugin.ManagerInterface
	objectStatus  ObjectStatus
}

func (o *objectNode) Create(ctx context.Context, object *unstructured.Unstructured) (*component.Node, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	apiVersion := object.GetAPIVersion()
	kind := object.GetKind()

	isReplicaSet, err := isObjectReplicaSet(object)
	if err != nil {
		return nil, err
	}

	if isReplicaSet {
		if err := checkReplicaCount(object); err != nil {
			return nil, err
		}
	}

	objectPath, err := o.objectPath(object)
	if err != nil {
		return nil, err
	}

	status, err := o.objectStatus.Status(ctx, object, o.link)
	if err != nil {
		return nil, err
	}

	node := &component.Node{
		Name:       accessor.GetName(),
		APIVersion: apiVersion,
		Kind:       kind,
		Status:     status.Status(),
		Properties: status.Properties,
		Details:    status.Details,
		Path:       objectPath,
	}

	return node, nil
}

func (o *objectNode) objectPath(object runtime.Object) (*component.Link, error) {
	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	return o.link.ForObjectWithQuery(object, accessor.GetName(), q)
}
