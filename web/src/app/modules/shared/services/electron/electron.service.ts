/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { ipcRenderer, webFrame } from 'electron';
import * as childProcess from 'child_process';
import * as fs from 'fs';

@Injectable({
  providedIn: 'root',
})
export class ElectronService {
  ipcRenderer: typeof ipcRenderer;
  webFrame: typeof webFrame;
  childProcess: typeof childProcess;
  fs: typeof fs;

  public portNumber: number;
  constructor() {
    if (this.isElectron()) {
      this.ipcRenderer = window.require('electron').ipcRenderer;
      this.webFrame = window.require('electron').webFrame;
      this.childProcess = window.require('child_process');
      this.fs = window.require('fs');

      this.ipcRenderer.once('port-message', (event, message) => {
        this.portNumber = message;
      });
    }
  }

  /**
   * Returns true if electron is detected
   */
  isElectron(): boolean {
    if (typeof process === 'undefined') {
      return false;
    }
    return (
      process && process.versions && process.versions.electron !== undefined
    );
  }

  /**
   * Returns the random port number from electron main process
   */
  port(): number {
    return this.portNumber;
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
