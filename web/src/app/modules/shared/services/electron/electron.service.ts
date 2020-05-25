/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ElectronService {
  constructor() {}

  /**
   * Returns true if electron is detected
   */
  isElectron(): boolean {
    return (
      process && process.versions && process.versions.electron !== undefined
    );
  }

  /**
   * Returns the platform.
   *   * Returns linux, darwin, or win32 for those platforms
   *   * Returns unknown if the platform is not linux, darwin, or win32
   *   * Returns a blank string is electron is not detected
   *
   */
  platform(): string {
    if (!this.isElectron()) {
      return '';
    }

    switch (process.platform) {
      case 'linux':
      case 'darwin':
      case 'win32':
        return process.platform;
      default:
        return 'unknown';
    }
  }
}
