/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/electron/preferences/v1alpha1"
	"github.com/vmware-tanzu/octant/internal/log"
)

// PreferencesUpdatedHandlerConfig is configuration for PreferencesUpdatedHandler.
type PreferencesUpdatedHandlerConfig interface {
	DataDirectory() string
}

// PreferencesUpdatedHandler is a handler for updated preferences.
type PreferencesUpdatedHandler struct {
	config PreferencesUpdatedHandlerConfig
}

var _ MessageHandler = &PreferencesUpdatedHandler{}

// NewPreferencesUpdatedHandler creates an instances of PreferencesUpdatedHandler.
func NewPreferencesUpdatedHandler(config PreferencesUpdatedHandlerConfig) *PreferencesUpdatedHandler {
	h := PreferencesUpdatedHandler{
		config: config,
	}

	return &h
}

// Key returns the key for this handler.
func (p PreferencesUpdatedHandler) Key() string {
	return "octant.cmd.updatePreferences"
}

// Handle handles an update request.
func (p PreferencesUpdatedHandler) Handle(ctx context.Context, in json.RawMessage) (interface{}, error) {
	logger := log.From(ctx)
	logger.Infof("updating preferences")

	var values map[string]string

	if err := json.Unmarshal(in, &values); err != nil {
		return nil, fmt.Errorf("unmarshal message: %w", err)
	}

	preferences, err := v1alpha1.CreateOrOpenPreferences(ctx, p.config)
	if err != nil {
		return nil, fmt.Errorf("create or open preferences: %w", err)
	}

	preferences.Update(values)

	if err := v1alpha1.WritePreferences(ctx, p.config, preferences); err != nil {
		return nil, fmt.Errorf("write preferences: %w", err)
	}

	return true, nil
}
