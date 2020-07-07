/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/internal/cluster"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

// ApplyYaml creates a yaml applier
type ApplyYaml struct {
	logger        log.Logger
	objectStore   store.Store
	clusterClient cluster.ClientInterface
}

var _ action.Dispatcher = (*ApplyYaml)(nil)

// NewApplyYaml creates an instance of ApplyYaml
func NewApplyYaml(logger log.Logger, objectStore store.Store, clusterClient cluster.ClientInterface) *ApplyYaml {
	return &ApplyYaml{
		logger:        logger,
		objectStore:   objectStore,
		clusterClient: clusterClient,
	}
}

// ActionName returns the name of this action
func (p *ApplyYaml) ActionName() string {
	return ActionApplyYaml
}

// Handle applies the requested yaml to the cluster
func (p *ApplyYaml) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	p.logger.With("payload", payload).Debugf("received action payload")
	request, err := applyYamlRequestFromPayload(payload)
	if err != nil {
		return errors.Wrap(err, "convert payload to apply yaml request")
	}
	p.logger.Debugf("%s", request)

	results := []string{}
	err = request.WithDoc(func(doc map[string]interface{}) error {
		p.logger.Debugf("apply resource %#v", doc)

		unstructuredObj := &unstructured.Unstructured{Object: doc}
		key, err := store.KeyFromObject(unstructuredObj)
		if err != nil {
			return err
		}
		gvr, namespaced, err := p.clusterClient.Resource(key.GroupVersionKind().GroupKind())
		if err != nil {
			return fmt.Errorf("unable to discover resource: %w", err)
		}
		if namespaced && key.Namespace == "" {
			unstructuredObj.SetNamespace(request.Namespace)
			key.Namespace = request.Namespace
		}

		if _, err := p.objectStore.Get(ctx, key); err != nil {
			if !kerrors.IsNotFound(err) {
				// unexpected error
				return fmt.Errorf("unable to get resource: %w", err)
			}

			// create object
			err := p.objectStore.Create(ctx, &unstructured.Unstructured{Object: doc})
			if err != nil {
				return fmt.Errorf("unable to create resource: %w", err)
			}

			result := fmt.Sprintf("Created %s (%s) %s", key.Kind, key.APIVersion, key.Name)
			if namespaced {
				result = fmt.Sprintf("%s in %s", result, key.Namespace)
			}
			results = append(results, result)

			return nil
		}

		// update object
		unstructuredYaml, err := sigyaml.Marshal(doc)
		if err != nil {
			return fmt.Errorf("unable to marshal resource as yaml: %w", err)
		}
		client, err := p.clusterClient.DynamicClient()
		if err != nil {
			return fmt.Errorf("unable to get dynamic client: %w", err)
		}

		withForce := true
		if namespaced {
			_, err = client.Resource(gvr).Namespace(key.Namespace).Patch(
				ctx,
				key.Name,
				types.ApplyPatchType,
				unstructuredYaml,
				metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
			)
			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}
		} else {
			_, err = client.Resource(gvr).Patch(
				ctx,
				key.Name,
				types.ApplyPatchType,
				unstructuredYaml,
				metav1.PatchOptions{FieldManager: "octant", Force: &withForce},
			)
			if err != nil {
				return fmt.Errorf("unable to patch resource: %w", err)
			}
		}

		result := fmt.Sprintf("Updated %s (%s) %s", key.Kind, key.APIVersion, key.Name)
		if namespaced {
			result = fmt.Sprintf("%s in %s", result, key.Namespace)
		}
		results = append(results, result)

		return nil
	})

	if err != nil {
		p.logger.Warnf("unable to apply yaml: %s", err)
		message := fmt.Sprintf("Unable to apply yaml: %s", err)
		alerter.SendAlert(action.CreateAlert(action.AlertTypeError, message, action.DefaultAlertExpiration))
		// do not return to send partial results to the client
	}

	switch len(results) {
	case 0:
		// nothing to do
	case 1:
		message := results[0]
		alerter.SendAlert(action.CreateAlert(action.AlertTypeInfo, message, action.DefaultAlertExpiration))
	default:
		// TODO(scothis) send detailed results to client
		message := fmt.Sprintf("Applied %d resources", len(results))
		alerter.SendAlert(action.CreateAlert(action.AlertTypeInfo, message, action.DefaultAlertExpiration))
	}
	return nil
}

type applyYamlRequest struct {
	Namespace string `json:"namespace,omitempty"`
	Update    string `json:"update,omitempty"`
}

func (req *applyYamlRequest) Validate() error {
	if req.Namespace == "" {
		return errors.New("namespace is blank")
	}

	if req.Update == "" {
		return errors.New("update is blank")
	}

	return nil
}

func (req *applyYamlRequest) WithDoc(cb func(doc map[string]interface{}) error) error {
	d := yaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(req.Update), 4096)
	for {
		doc := map[string]interface{}{}
		if err := d.Decode(&doc); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("unable to parse yaml: %w", err)
		}
		if len(doc) == 0 {
			// skip empty documents
			continue
		}
		if err := cb(doc); err != nil {
			return err
		}
	}
}

func applyYamlRequestFromPayload(payload action.Payload) (*applyYamlRequest, error) {
	namespace, err := payload.String("namespace")
	if err != nil {
		return nil, err
	}

	update, err := payload.String("update")
	if err != nil {
		return nil, err
	}

	req := &applyYamlRequest{
		Namespace: namespace,
		Update:    update,
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return req, nil
}
