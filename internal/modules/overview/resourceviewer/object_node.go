package resourceviewer

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/pkg/plugin"
	"github.com/vmware/octant/pkg/view/component"
)

type objectNode struct {
	link          link.Interface
	pluginPrinter plugin.ManagerInterface
	objectStatus  ObjectStatus
}

func (o *objectNode) Create(ctx context.Context, object runtime.Object) (*component.Node, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	groupVersionKind := object.GetObjectKind().GroupVersionKind()
	apiVersion, kind := groupVersionKind.ToAPIVersionAndKind()

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

	status, err := o.objectStatus.Status(ctx, object)
	if err != nil {
		return nil, err
	}

	node := &component.Node{
		Name:       accessor.GetName(),
		APIVersion: apiVersion,
		Kind:       kind,
		Status:     status.Status(),
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
