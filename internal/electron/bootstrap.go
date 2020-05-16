/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"context"
	"fmt"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"

	"github.com/vmware-tanzu/octant/internal/electron/preferences/v1alpha1"
	"github.com/vmware-tanzu/octant/internal/log"
)

func dimensionsWindowOptions(in astilectron.WindowOptions) astilectron.WindowOptions {
	height := 750
	width := 1200
	windowMinWidth := 768
	windowMinHeight := windowMinWidth * 10 / 16 // base min height off ultra wide ratio

	in.Height = astikit.IntPtr(height)
	in.Width = astikit.IntPtr(width)
	in.MinWidth = astikit.IntPtr(windowMinWidth)
	in.MinHeight = astikit.IntPtr(windowMinHeight)

	return in
}

func initWindows(ctx context.Context, a *astilectron.Astilectron, appURL string, listener MessageListener, logger astikit.SeverityLogger) ([]*astilectron.Window, error) {
	windowOptions := astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
	}

	windowOptions = dimensionsWindowOptions(windowOptions)
	windowOptions = platformWindowOptions(windowOptions)

	window, err := a.NewWindow(appURL, &windowOptions)

	if err != nil {
		return nil, fmt.Errorf("create main window: %w", err)
	}

	if listener != nil {
		window.OnMessage(handleMessage(ctx, window, listener, logger))
	}

	windows := []*astilectron.Window{
		window,
	}

	return windows, nil
}

func initMenuItems(ctx context.Context, window *astilectron.Window, paths astilectron.Paths) ([]*astilectron.MenuItemOptions, error) {
	logger := log.From(ctx)

	menuItems := []*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About Octant"),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Preferences"),
					OnClick: func(e astilectron.Event) bool {
						logger.Infof("open preferences event")
						preferences, err := v1alpha1.CreateOrOpenPreferences(ctx, paths)
						if err != nil {
							logger.WithErr(err).Errorf("create or open preferences")
							return false
						}

						SendMessage(ctx, window, OctantCmdPreferences, preferences)
						return false
					},
					Accelerator: astilectron.NewAccelerator("CommandOrControl", ","),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Services"),
					Role:  astilectron.MenuItemRoleServices,
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Hide Octant"),
					Role:  astilectron.MenuItemRoleHide,
				},
				{
					Label: astikit.StrPtr("Hide Others"),
					Role:  astilectron.MenuItemRoleHideOthers,
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Quit"),
					Role:  astilectron.MenuItemRoleQuit,
				},
			},
		},
		{
			Label: astikit.StrPtr("View"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Log"),
				},
				{
					Type: astilectron.MenuItemTypeSeparator,
				},
				{
					Label: astikit.StrPtr("Develop"),
					SubMenu: []*astilectron.MenuItemOptions{
						{
							Label:       astikit.StrPtr("Developer Tools"),
							Accelerator: astilectron.NewAccelerator("CommandOrControl", "Option", "I"),
							OnClick: func(e astilectron.Event) bool {
								if err := window.OpenDevTools(); err != nil {
									logger.WithErr(err).Errorf("open dev tools")
									return false
								}

								return false
							},
						},
					},
				},
			},
		},

		{
			Label: astikit.StrPtr("Window"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("Minimize"),
					Role:  astilectron.MenuItemRoleMinimize,
				},
				{
					Label: astikit.StrPtr("Zoom"),
					Role:  astilectron.MenuItemRoleZoom,
				},
			},
		},
	}

	return menuItems, nil
}
