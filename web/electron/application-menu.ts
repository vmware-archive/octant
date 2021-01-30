/*
 *  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app, Menu, shell } from 'electron';
import { errLogPath } from './paths';

export class ApplicationMenu {
  public readonly menu: Menu;

  constructor() {
    this.menu = this.createMenu();
  }

  createMenu() {
    const template: Electron.MenuItemConstructorOptions[] = [
      {
        label: 'File',
        submenu: [
          {
            label: 'Close',
            role: 'minimize',
            accelerator: 'CommandOrControl+w',
          },
          {
            label: 'Quit Octant',
            accelerator: 'CommandOrControl+q',
            click() {
              app.quit();
            },
          },
        ],
      },
      {
        label: 'Edit',
        submenu: [
          { role: 'undo' },
          { role: 'redo' },
          { role: 'cut' },
          { role: 'copy' },
          { role: 'paste' },
        ],
      },
      {
        label: 'View',
        submenu: [
          { role: 'resetZoom' },
          { role: 'zoomIn', accelerator: 'CommandOrControl+=' },
          { role: 'zoomOut' },
          { type: 'separator' },
          { role: 'togglefullscreen', accelerator: 'CommandOrControl+Shift+F' },
          { role: 'toggleDevTools' },
          { type: 'separator' },
          {
            label: 'View Logs',
            click() {
              shell.showItemInFolder(errLogPath);
            },
          },
        ],
      },
      {
        label: 'Help',
        submenu: [
          {
            label: 'octant.dev',
            click() {
              shell.openExternal('https://octant.dev/');
            },
          },
        ],
      },
    ];

    const menu = Menu.buildFromTemplate(template);
    return menu;
  }
}
