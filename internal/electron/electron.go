/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	astibundler "github.com/asticode/go-astilectron-bundler"
)

const defaultResourcePath = "resources"

// Options are options for configuring the electron app.
type Options struct {
	// AppName is the name of the application.
	AppName string
	// Asset is asset loader.
	Asset AssetFunc
	// AssetDir is the asset directory lister.
	AssetDir AssetDirFunc
	// RestoreAssets is the asset restorer.
	RestoreAssets RestoreAssetsFunc
	// VersionAstilectron is the astilectron version.
	VersionAstilectron string
	// VersionElectron is the electron version.
	VersionElectron string
}

// Electron manages the electron app.
type Electron struct {
	instance *astilectron.Astilectron
	windows  []*astilectron.Window
	menu     *astilectron.Menu
	listener MessageListener
}

// New creates an instance of Electron.
func New(ctx context.Context, options Options) (*Electron, error) {
	a, err := initAstilectron(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("initialize electron: %w", err)
	}

	listener := NewMessageListener()

	e := Electron{
		instance: a,
		listener: listener,
	}

	listener.Register(NewPreferencesUpdatedHandler(a.Paths()))

	return &e, nil
}

// Start starts the electron app.
func (e *Electron) Start(ctx context.Context, appURL string) error {
	logger := loggerFromContext(ctx)

	instance := e.instance

	if err := instance.Start(); err != nil {
		return fmt.Errorf("starting electron failed: %w", err)
	}

	// init windows
	windows, err := initWindows(ctx, instance, appURL, e.listener, logger)
	if err != nil {
		return fmt.Errorf("build windows")
	}

	for i := range windows {
		if err := windows[i].Create(); err != nil {
			return fmt.Errorf("create window %d: %w", i, err)
		}
	}

	e.windows = windows

	// create menu options
	menuItems, err := initMenuItems(ctx, windows[0], e.instance.Paths())
	if err != nil {
		return fmt.Errorf("init menu items: %w", err)
	}

	// setup menu
	menu := instance.NewMenu(menuItems)
	if err := menu.Create(); err != nil {
		return fmt.Errorf("create menu: %w", err)
	}

	e.menu = menu

	// setup tray

	return nil
}

// Stop stops the electron app.
func (e *Electron) Stop() {
	e.instance.Stop()
}

// Wait waits for the electron app to stop.
func (e *Electron) Wait() {
	e.instance.Wait()
}

func initAstilectron(ctx context.Context, options Options) (*astilectron.Astilectron, error) {
	logger := loggerFromContext(ctx)

	aOptions := astilectron.Options{
		AppName:            options.AppName,
		AppIconDarwinPath:  filepath.Join(defaultResourcePath, "icon.icns"),
		AppIconDefaultPath: filepath.Join(defaultResourcePath, "icon.png"),
		SingleInstance:     true,
		VersionAstilectron: options.VersionAstilectron,
		VersionElectron:    options.VersionElectron,
	}

	a, err := astilectron.New(logger, aOptions)
	if err != nil {
		return nil, fmt.Errorf("create astilectron: %w", err)
	}

	a.HandleSignals(astikit.LoggerSignalHandler(logger))

	if options.Asset != nil {
		a.SetProvisioner(astibundler.NewProvisioner(options.Asset, logger))
	}

	if options.RestoreAssets != nil {
		if err := restoreResources(logger, a, options.Asset, options.AssetDir, options.RestoreAssets, defaultResourcePath); err != nil {
			return nil, fmt.Errorf("restore resources: %w", err)
		}
	}

	return a, nil
}
