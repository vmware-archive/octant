/*
*  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
*  SPDX-License-Identifier: Apache-2.0
*
*/

import ElectronStore = require('electron-store');

interface OctantStore {
  minimizeToTray: boolean;
  showDialogue: boolean;
  theme: string;
  navigation: {
    collapsed: boolean;
    labels: boolean;
  }
  windowBounds: Electron.Rectangle
}


export const electronStore = new ElectronStore<OctantStore>({
  defaults: {
    minimizeToTray: true,
    showDialogue: true,
    theme: 'light',
    windowBounds: undefined,
    navigation: {
      collapsed: false,
      labels: true,
    },
  }
});

