/*
*  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
*  SPDX-License-Identifier: Apache-2.0
*
*/

import ElectronStore = require('electron-store');

interface OctantStore {
  minimizeToTray: boolean;
  showDialogue: boolean;
}


export const electronStore = new ElectronStore<OctantStore>({
  defaults: {
    minimizeToTray: true,
    showDialogue: true
  }
});

//
// export class OctantConfig {
//   constructor() {
//   }
//
//   set() {}
//   get() {}
// }
