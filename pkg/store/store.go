/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"

	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/cluster"
)

//go:generate mockgen  -destination=./fake/mock_store.go -package=fake github.com/vmware-tanzu/octant/pkg/store Store

// UpdateFn is a function that is called when
type UpdateFn func(store Store)

// Store stores Kubernetes objects.
type Store interface {
	List(ctx context.Context, key Key) (list *unstructured.UnstructuredList, loading bool, err error)
	Get(ctx context.Context, key Key) (object *unstructured.Unstructured, err error)
	Delete(ctx context.Context, key Key) error
	Watch(ctx context.Context, key Key, handler cache.ResourceEventHandler) error
	Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error
	UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error
	Update(ctx context.Context, key Key, updater func(*unstructured.Unstructured) error) error
	IsLoading(ctx context.Context, key Key) bool
	Create(ctx context.Context, object *unstructured.Unstructured) error
	// CreateOrUpdateFromYAML creates resources in the cluster from YAML input.
	// Resources are created in the order they are present in the YAML.
	// An error creating a resource halts resource creation.
	// A list of created resources is returned. You may have created resources AND a non-nil error.
	CreateOrUpdateFromYAML(ctx context.Context, namespace, input string) ([]string, error)
}

// Key is a key for the object store.
type Key struct {
	Namespace     string                `json:"namespace"`
	APIVersion    string                `json:"apiVersion"`
	Kind          string                `json:"kind"`
	Name          string                `json:"name"`
	Selector      *labels.Set           `json:"selector"`
	LabelSelector *metav1.LabelSelector `json:"labelSelector"`
}

// Validate validates the key.
func (k Key) Validate() error {
	var err error

	if k.APIVersion == "" {
		err = multierror.Append(err, errors.New("apiVersion is blank"))
	}

	if k.Kind == "" {
		err = multierror.Append(err, errors.New("kind is blank"))
	}

	if k.LabelSelector != nil && len(k.LabelSelector.MatchExpressions) > 0 {
		for _, v := range k.LabelSelector.MatchExpressions {
			if (v.Operator == metav1.LabelSelectorOpIn || v.Operator == metav1.LabelSelectorOpNotIn) && len(v.Values) == 0 {
				err = multierror.Append(err, errors.New("operator In/NotIn must not have empty values array"))
			} else if (v.Operator == metav1.LabelSelectorOpExists || v.Operator == metav1.LabelSelectorOpDoesNotExist) && len(v.Values) > 0 {
				err = multierror.Append(err, errors.New("operator Exists/DoesNotExist must have empty values array"))
			}
		}
	}

	return err
}

func (k Key) String() string {
	var sb strings.Builder

	sb.WriteString("CacheKey[")
	if k.Namespace != "" {
		sb.WriteString(fmt.Sprintf("Namespace='%s', ", k.Namespace))
	}
	sb.WriteString(fmt.Sprintf("APIVersion='%s', ", k.APIVersion))
	sb.WriteString(fmt.Sprintf("Kind='%s'", k.Kind))

	if k.Name != "" {
		sb.WriteString(fmt.Sprintf(", Name='%s'", k.Name))
	}

	if k.Selector != nil && k.Selector.String() != "" {
		sb.WriteString(fmt.Sprintf(", Selector='%s'", k.Selector.String()))
	}

	if k.LabelSelector != nil {
		sb.WriteString(", LabelSelector='")
		k.labelSelectorString(&sb)
		sb.WriteString("'")
	}

	sb.WriteString("]")

	return sb.String()
}

func (k Key) labelSelectorString(sb *strings.Builder) {
	if k.LabelSelector.MatchLabels != nil {
		sb.WriteString("MatchLabels=")
		data := make([]string, len(k.LabelSelector.MatchLabels))
		for k, v := range k.LabelSelector.MatchLabels {
			data = append(data, fmt.Sprintf("%s=%s", k, v))
		}
		sb.WriteString(strings.Join(data, ","))
	}

	if k.LabelSelector.MatchLabels != nil && k.LabelSelector.MatchExpressions != nil {
		sb.WriteString(",")
	}

	if k.LabelSelector.MatchExpressions != nil {
		sb.WriteString("MatchExpressions=")
		data := make([]string, len(k.LabelSelector.MatchLabels))
		for _, me := range k.LabelSelector.MatchExpressions {
			data = append(data, fmt.Sprintf("%s %s (%s)", me.Key, me.Operator, strings.Join(me.Values, ",")))
		}
		sb.WriteString(strings.Join(data, ","))
	}
}

// GroupVersionKind converts the Key to a GroupVersionKind.
func (k Key) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(k.APIVersion, k.Kind)
}

// ToActionPayload converts the Key to a payload.
func (k Key) ToActionPayload() action.Payload {
	return action.Payload{
		"namespace":  k.Namespace,
		"apiVersion": k.APIVersion,
		"kind":       k.Kind,
		"name":       k.Name,
	}
}

// KeyFromPayload converts a payload into a Key.
func KeyFromPayload(payload action.Payload) (Key, error) {
	namespace, err := payload.OptionalString("namespace")
	if err != nil {
		return Key{}, err
	}
	apiVersion, err := payload.String("apiVersion")
	if err != nil {
		return Key{}, err
	}
	kind, err := payload.String("kind")
	if err != nil {
		return Key{}, err
	}
	name, err := payload.OptionalString("name")
	if err != nil {
		return Key{}, err
	}

	key := Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	labelSelectorBytes, err := payload.Raw("labelSelector")
	if err == nil {
		labelSelector := metav1.LabelSelector{}
		if err := json.Unmarshal(labelSelectorBytes, &labelSelector); err != nil {
			return Key{}, fmt.Errorf("label selector contents are invalid: %w", err)
		}

		key.LabelSelector = &labelSelector
	}

	labelSetBytes, err := payload.Raw("selector")
	if err == nil {
		set := labels.Set{}
		if err := json.Unmarshal(labelSetBytes, &set); err != nil {
			return Key{}, fmt.Errorf("selector contents are invalid: %w", err)
		}

		key.Selector = &set
	}

	return key, nil
}

// KeyFromObject creates a key from a runtime object.
func KeyFromObject(object runtime.Object) (Key, error) {
	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return Key{}, err
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return Key{}, err
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return Key{}, err
	}

	name, err := accessor.Name(object)
	if err != nil {
		return Key{}, err
	}

	return Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}, nil
}

// KeyFromGroupVersionKind creates a key from a group version kind.
func KeyFromGroupVersionKind(groupVersionKind schema.GroupVersionKind) Key {
	apiVersion, kind := groupVersionKind.ToAPIVersionAndKind()

	return Key{
		APIVersion: apiVersion,
		Kind:       kind,
	}
}
