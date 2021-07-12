/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

package dash

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fsnotify/fsnotify"

	"github.com/vmware-tanzu/octant/internal/log"
	internalStrings "github.com/vmware-tanzu/octant/internal/util/strings"
)

//go:generate mockgen -destination=./fake/mock_file_watcher.go -package=fake github.com/vmware-tanzu/octant/pkg/dash FileWatcher
//go:generate mockgen -destination=./fake/mock_watcher_config.go -package=fake github.com/vmware-tanzu/octant/pkg/dash WatcherConfig

// FileWatcher watches files.
type FileWatcher interface {
	// Add adds a file to watch.
	Add(name string) error
	// Events returns a channel with fsnotify events.
	Events() chan fsnotify.Event
	// Errors returns errors.
	Errors() chan error
}

// DefaultFileWatcher wraps fsnotify.Watcher so it can adhere to the
// FileWatcher API.
type DefaultFileWatcher struct {
	watcher *fsnotify.Watcher
}

var _ FileWatcher = (*DefaultFileWatcher)(nil)

// NewDefaultFileWatcher creates an instance of DefaultFileWatcher.
func NewDefaultFileWatcher() (*DefaultFileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	return &DefaultFileWatcher{
		watcher: watcher,
	}, nil
}

// Add adds a file name to be watched.
func (d *DefaultFileWatcher) Add(name string) error {
	return d.watcher.Add(name)
}

// Events returns a channel of fsnotify events.
func (d *DefaultFileWatcher) Events() chan fsnotify.Event {
	return d.watcher.Events
}

// Errors returns a channel of errors.
func (d DefaultFileWatcher) Errors() chan error {
	return d.watcher.Errors
}

// ConfigWatcherOption is an option for configuration ConfigWatcher.
type ConfigWatcherOption func(cw *ConfigWatcher)

// ConfigWatcherFileWatcher sets the file watcher for ConfigWatcher.
func ConfigWatcherFileWatcher(fw FileWatcher) ConfigWatcherOption {
	return func(cw *ConfigWatcher) {
		cw.FileWatcher = fw
	}
}

// WatcherConfig is an interface with configuration for ConfigWatcher.
type WatcherConfig interface {
	CurrentContext() string
	UseFSContext(ctx context.Context) error
	UseContext(ctx context.Context, name string) error
}

// ConfigWatcher watches kubernetes configurations.
type ConfigWatcher struct {
	FileWatcher   FileWatcher
	watcherConfig WatcherConfig
}

// NewConfigWatcher creates an instance of ConfigWatcher.
func NewConfigWatcher(wc WatcherConfig, options ...ConfigWatcherOption) (*ConfigWatcher, error) {
	cw := &ConfigWatcher{
		watcherConfig: wc,
	}

	for _, option := range options {
		option(cw)
	}

	if cw.FileWatcher == nil {
		fw, err := NewDefaultFileWatcher()
		if err != nil {
			return nil, fmt.Errorf("create file watcher: %w", err)
		}
		cw.FileWatcher = fw
	}

	return cw, nil
}

// Add adds file names to be watched.
func (cw *ConfigWatcher) Add(ctx context.Context, names ...string) error {
	logger := log.From(ctx).With("component", "config-watcher")

	for _, name := range names {
		logger.With("config", name).Infof("watching config file")
		if err := cw.FileWatcher.Add(name); err != nil {
			return fmt.Errorf("unable to watch %s: %w", name, err)
		}
	}

	return nil
}

// Watch runs the config watcher loop
func (cw *ConfigWatcher) Watch(ctx context.Context) {
	logger := log.From(ctx).With("component", "config-watcher")
	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
			logger.Infof("shutting down config watcher")
		case <-cw.FileWatcher.Events():
			if err := cw.watcherConfig.UseFSContext(ctx); err != nil {
				logger.WithErr(err).Errorf("reload config")
			}
		case err := <-cw.FileWatcher.Errors():
			logger.WithErr(err).Errorf("event error")
		}
	}
}

// watchConfigs watches kubernetes config files. If any file changes, reload the cluster client.
func watchConfigs(ctx context.Context, wc WatcherConfig, configChainStr string) error {
	cw, err := NewConfigWatcher(wc)
	if err != nil {
		return fmt.Errorf("create config watcher: %w", err)
	}

	chain := internalStrings.Deduplicate(filepath.SplitList(configChainStr))
	for _, f := range chain {
		if err := cw.Add(ctx, f); err != nil {
			return fmt.Errorf("unable to watch %s: %w", f, err)
		}
	}

	go cw.Watch(ctx)

	return nil
}
