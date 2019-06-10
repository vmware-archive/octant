package link

import (
	"net/url"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/heptio/developer-dash/pkg/view/component"
)

//go:generate mockgen -source=link.go -destination=./fake/mock_link.go -package=fake github.com/heptio/developer-dash/internal/overview/link

type objectPathFn func(namespace, apiVersion, kind, name string) (string, error)

type Interface interface {
	ForObject(object runtime.Object, text string) (*component.Link, error)
	ForObjectWithQuery(object runtime.Object, text string, query url.Values) (*component.Link, error)
	ForGVK(namespace, apiVersion, kind, name, text string) (*component.Link, error)
	ForOwner(parent runtime.Object, controllerRef *metav1.OwnerReference) (*component.Link, error)
}

type Config interface {
	ObjectPath(namespace, apiVersion, kind, name string) (string, error)
}

type Link struct {
	objectPathFn objectPathFn
}

var _ Interface = (*Link)(nil)

func NewFromDashConfig(config Config) (*Link, error) {
	if config == nil {
		return nil, errors.New("link config is nil")
	}

	return &Link{
		objectPathFn: config.ObjectPath,
	}, nil
}

// ForObject returns a link component referencing an object
// Returns an empty link if an error occurs.
func (l *Link) ForObject(object runtime.Object, text string) (*component.Link, error) {
	p, err := l.extractPathFromObject(object)
	if err != nil {
		return nil, err
	}

	return component.NewLink("", text, p), nil
}

// ForObjectWithQuery returns a link component references an object with a query.
// Return an empty link if an error occurs.
func (l *Link) ForObjectWithQuery(object runtime.Object, text string, query url.Values) (*component.Link, error) {
	p, err := l.extractPathFromObject(object)
	if err != nil {
		return nil, err
	}

	u := url.URL{Path: p}
	u.RawQuery = query.Encode()
	return component.NewLink("", text, u.String()), nil
}

// ForGVK returns a link component referencing an object
func (l *Link) ForGVK(namespace, apiVersion, kind, name, text string) (*component.Link, error) {
	p, err := l.objectPathFn(namespace, apiVersion, kind, name)
	if err != nil {
		return nil, err
	}

	return component.NewLink("", text, p), nil
}

// ForOwner returns a link component for an owner.
func (l *Link) ForOwner(parent runtime.Object, controllerRef *metav1.OwnerReference) (*component.Link, error) {
	if controllerRef == nil || parent == nil {
		return component.NewLink("", "none", ""), nil
	}

	accessor := meta.NewAccessor()
	ns, err := accessor.Namespace(parent)
	if err != nil {
		return component.NewLink("", "none", ""), nil
	}

	return l.ForGVK(
		ns,
		controllerRef.APIVersion,
		controllerRef.Kind,
		controllerRef.Name,
		controllerRef.Name,
	)
}

func (l *Link) extractPathFromObject(object runtime.Object) (string, error) {
	if object == nil {
		return "", errors.New("can't generate path for nil object")
	}

	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return "", err
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return "", err
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return "", err
	}

	name, err := accessor.Name(object)
	if err != nil {
		return "", err
	}

	return l.objectPathFn(namespace, apiVersion, kind, name)
}
