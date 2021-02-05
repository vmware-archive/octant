/*
 *  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app, BrowserWindow, Menu, Tray, nativeImage, shell } from 'electron';
import { errLogPath, iconPath } from './paths';

export class TrayMenu {
  public readonly tray: Tray;

  constructor(public window: BrowserWindow) {
    this.tray = new Tray(this.createNativeImage());
    this.tray.setContextMenu(this.createMenu(window));
  }

  createNativeImage(): Electron.NativeImage {
    const image = nativeImage.createFromPath(iconPath);
    image.setTemplateImage(true);
    return image.resize({ width: 16, height: 16 });
  }

  createMenu(win: BrowserWindow): Menu {
    const menu = Menu.buildFromTemplate([
      {
        label: 'Open Octant',
        type: 'normal',
        click: () => {
          win.show();
        },
      },
      {
        label: 'View Logs',
        type: 'normal',
        click: () => {
          shell.showItemInFolder(errLogPath);
        },
      },
      {
        label: 'Quit',
        type: 'normal',
        click: () => {
          app.quit();
        },
      },
    ]);
    return menu;
  }
}
