/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package octant

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

// ApplyYaml creates a yaml applier
type ApplyYaml struct {
	logger      log.Logger
	objectStore store.Store
}

var _ action.Dispatcher = (*ApplyYaml)(nil)

// NewApplyYaml creates an instance of ApplyYaml
func NewApplyYaml(logger log.Logger, objectStore store.Store) *ApplyYaml {
	return &ApplyYaml{
		logger:      logger,
		objectStore: objectStore,
	}
}

// ActionName returns the name of this action
func (p *ApplyYaml) ActionName() string {
	return action.ActionApplyYaml
}

// Handle applies the requested yaml to the cluster
func (p *ApplyYaml) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	p.logger.With("payload", payload).Debugf("received action payload")

	request, err := applyYamlRequestFromPayload(payload)
	if err != nil {
		return errors.Wrap(err, "convert payload to apply yaml request")
	}
	p.logger.Debugf("%s", request)

	results, err := p.objectStore.CreateOrUpdateFromYAML(ctx, request.Namespace, request.Update)
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
