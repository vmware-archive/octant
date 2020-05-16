/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/octant/internal/log"
)

const (
	// PreferencesName is the name for the preferences file.
	PreferencesName = "preferences.json"
)

// Config is configuration for preferences.
type Config interface {
	// DataDirectory is the data directory.
	DataDirectory() string
}

// WritePreferences writes preferences to the filesystem.
func WritePreferences(ctx context.Context, config Config, preferences Preferences) error {
	logger := log.From(ctx)
	logger = logger.With("preferencesFilename", preferencesFileName(config))

	logger.Infof("writing preferences")

	preferencesFilename := filepath.Join(config.DataDirectory(), PreferencesName)

	f, err := os.OpenFile(preferencesFilename, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open preferences file for writing: %w", err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			logger.WithErr(cErr).Errorf("close preferences file")
		}
	}()

	if err := json.NewEncoder(f).Encode(&preferences); err != nil {
		return fmt.Errorf("encode preferences: %w", err)
	}

	return nil
}

// CreateOrOpenPreferences creates or opens an existing preferences file.
func CreateOrOpenPreferences(ctx context.Context, config Config) (Preferences, error) {
	logger := log.From(ctx)

	logger = logger.With("preferencesFilename", preferencesFileName(config))

	logger.Infof("checking preferences")
	fi, err := os.Stat(preferencesFileName(config))
	if err != nil {
		if !os.IsNotExist(err) {
			return Preferences{}, fmt.Errorf("check if preferences exists: %w", err)
		}

		logger.Infof("initializing preferences")
		if err := initPreferences(ctx, config); err != nil {
			return Preferences{}, fmt.Errorf("initialize preferences: %w", err)
		}

	}

	if fi != nil && fi.IsDir() {
		return Preferences{}, fmt.Errorf("preferences file is a directory")
	}

	logger.Infof("open preferences")
	f, err := os.Open(preferencesFileName(config))
	if err != nil {
		return Preferences{}, fmt.Errorf("open preferences file: %w", err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			logger.WithErr(cErr).Errorf("close preferences file")
		}
	}()

	var p Preferences
	if err := json.NewDecoder(f).Decode(&p); err != nil {
		return Preferences{}, fmt.Errorf("decode preferences: %w", err)
	}

	return p, nil
}

func preferencesFileName(config Config) string {
	return filepath.Join(config.DataDirectory(), PreferencesName)
}

func initPreferences(ctx context.Context, config Config) error {
	logger := log.From(ctx)
	logger.With("filename", preferencesFileName(config)).Infof("creating preference file")

	if err := WritePreferences(ctx, config, defaultPreferences()); err != nil {
		return fmt.Errorf("write preferences: %w", err)
	}

	return nil
}

func defaultPreferences() Preferences {
	embeddedChoices := map[string]string{
		"Embedded": "embedded",
		"Proxied":  "proxied",
	}

	proxyURLConditions := []Condition{
		{
			LHS: "development.embedded",
			RHS: "proxied",
			Op:  OperationTypeString,
		},
	}

	p := Preferences{
		Version:    "v1alpha1",
		UpdateName: "octant.cmd.updatePreferences",
		Panels: []PreferencePanel{
			{
				Name: "Development",
				Sections: []PreferenceSection{
					{
						Name: "Frontend Source",
						Elements: []Element{
							NewRadioElement("development.embedded", "embedded", nil,
								embeddedChoices),
							NewTextElement("development.frontendProxyURL", "", proxyURLConditions,
								"http://localhost:4200", "Frontend Proxy URL"),
						},
					},
				},
			},
		},
	}

	return p
}
