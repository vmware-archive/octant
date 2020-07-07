/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstore

import (
	"context"
	"fmt"
	"sync"

	"go.opencensus.io/trace"
	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/cluster"
	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/pkg/store"
)

type AccessError struct {
	Key AccessKey
}

func (ae *AccessError) Error() string {
	return fmt.Sprintf("access denied: no %s access in %s to %s/%s",
		ae.Key.Verb, ae.Key.Namespace, ae.Key.Group, ae.Key.Resource)
}

// AccessKey is used at a key in an access map. It is made up of a Namespace, Group, Resource, and Verb.
type AccessKey struct {
	Namespace string
	Group     string
	Resource  string
	Verb      string
}

type accessMap map[AccessKey]bool

type accessCache struct {
	access accessMap
	mu     sync.RWMutex
}

func newAccessCache() *accessCache {
	return &accessCache{
		access: accessMap{},
	}
}

func (c *accessCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.access = accessMap{}
}

func (c *accessCache) set(key AccessKey, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.access[key] = value
}

func (c *accessCache) get(key AccessKey) (v, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok = c.access[key]
	return v, ok
}

type ResourceAccess interface {
	HasAccess(context.Context, store.Key, string) error
	Reset()
	Get(AccessKey) (bool, bool)
	Set(AccessKey, bool)
	UpdateClient(client cluster.ClientInterface)
}

type resourceAccess struct {
	client cluster.ClientInterface
	cache  *accessCache

	mu sync.RWMutex
}

var _ ResourceAccess = (*resourceAccess)(nil)

func NewResourceAccess(client cluster.ClientInterface) ResourceAccess {
	return &resourceAccess{
		client: client,
		cache:  newAccessCache(),
	}
}

func (r *resourceAccess) UpdateClient(client cluster.ClientInterface) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.client = client
}

// Reset resets the resource access cache.
func (r *resourceAccess) Reset() {
	r.cache.reset()
}

// Get returns the value and if it was found for an AccessKey.
func (r *resourceAccess) Get(key AccessKey) (value, found bool) {
	return r.cache.get(key)
}

// Set will set the value in the map for an AccessKey.
func (r *resourceAccess) Set(key AccessKey, v bool) {
	r.cache.set(key, v)
}

// HasAccess returns an error if the current user does not have access to perform the verb action
// for the given key.
func (r *resourceAccess) HasAccess(ctx context.Context, key store.Key, verb string) error {
	_, span := trace.StartSpan(ctx, "resourceAccessHasAccess")
	defer span.End()

	aKey, err := r.keyToAccessKey(key, verb)
	if err != nil {
		return err
	}

	access, ok := r.cache.get(aKey)

	if !ok {
		span.Annotate([]trace.Attribute{}, "fetch access start")
		val, err := r.fetchAccess(aKey, verb)
		if err != nil {
			return fmt.Errorf("fetch access: %+v: %w", aKey, err)
		}

		r.cache.set(aKey, val)
		access = val
		span.Annotate([]trace.Attribute{}, "fetch access finish")
	}

	if !access {
		return oerrors.NewAccessError(key, verb, err)
	}

	return nil
}

func (r *resourceAccess) keyToAccessKey(key store.Key, verb string) (AccessKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gvk := key.GroupVersionKind()

	if gvk.GroupKind().Empty() {
		return AccessKey{}, fmt.Errorf("unable to check access for key %s", key.String())
	}

	gvr, _, err := r.client.Resource(gvk.GroupKind())
	if err != nil {
		return AccessKey{}, fmt.Errorf("unable to get resource for group kind %s: %w", gvk.GroupKind(), err)
	}

	aKey := AccessKey{
		Namespace: key.Namespace,
		Group:     gvr.Group,
		Resource:  gvr.Resource,
		Verb:      verb,
	}
	return aKey, nil
}

func (r *resourceAccess) fetchAccess(key AccessKey, verb string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	k8sClient, err := r.client.KubernetesClient()
	if err != nil {
		return false, fmt.Errorf("client kubernetes: %w", err)
	}

	authClient := k8sClient.AuthorizationV1()
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: key.Namespace,
				Group:     key.Group,
				Resource:  key.Resource,
				Verb:      verb,
			},
		},
	}

	review, err := authClient.SelfSubjectAccessReviews().Create(context.TODO(), sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("client auth: %w", err)
	}
	return review.Status.Allowed, nil
}
