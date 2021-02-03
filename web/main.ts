/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import {
  app,
  BrowserWindow,
  dialog,
  Menu,
  MessageBoxOptions,
  screen,
  session,
  shell,
} from 'electron';
import { ApplicationMenu } from './electron/application-menu';
import { TrayMenu } from './electron/tray-menu';
import { apiLogPath, errLogPath, tmpPath, iconPath } from './electron/paths';
import { electronStore } from './electron/store';
import * as path from 'path';
import * as child_process from 'child_process';
import * as process from 'process';
import * as os from 'os';
import * as fs from 'fs';

let win: BrowserWindow = null;
let serverPid: any = null;
let state: any = {};
let statePath: string = null;
let closing = false;
let tray = null;

const args = process.argv.slice(1);
const local = args.some(val => val === '--local');

const applicationMenu = new ApplicationMenu();
Menu.setApplicationMenu(applicationMenu.menu);

let saveBoundsCookie;

function saveCurrentState() {
  fs.writeFileSync(statePath, JSON.stringify(state, null, 2));
}

function loadSavedState() {
  try {
    state = JSON.parse(fs.readFileSync(statePath, 'utf8'));
  } catch (e) {}
}

function saveBoundsSoon() {
  if (saveBoundsCookie) clearTimeout(saveBoundsCookie);
  saveBoundsCookie = setTimeout(() => {
    saveBoundsCookie = undefined;
    state.bounds = win.getBounds();
    saveCurrentState();
  }, 1000);
}

function createWindow(): BrowserWindow {
  const electronScreen = screen;
  const size = electronScreen.getPrimaryDisplay().workAreaSize;

  const options = {
    x: null,
    y: null,
    width: size.width,
    height: size.height,
    title: '',
    webPreferences: {
      nodeIntegration: true,
      webSecurity: false,
      allowRunningInsecureContent: true,
      contextIsolation: false, // false if you want to run 2e2 test with Spectron
      enableRemoteModule: true, // true if you want to run 2e2 test  with Spectron or use remote module in renderer context (ie. Angular)
      // preload: path.join(__dirname, 'preload.js'),
    },
  };

  if (state?.bounds) {
    const bounds = state.bounds;
    const area = electronScreen.getDisplayMatching(bounds).workArea;
    if (
      bounds.x >= area.x &&
      bounds.y >= area.y &&
      bounds.x + bounds.width <= area.x + area.width &&
      bounds.y + bounds.height <= area.y + area.height
    ) {
      options.x = bounds.x;
      options.y = bounds.y;
    }
    // If the saved size is still valid, use it.
    if (bounds.width <= area.width || bounds.height <= area.height) {
      options.width = bounds.width;
      options.height = bounds.height;
    }
  }

  // Create the browser window.
  win = new BrowserWindow(options);
  win.setIcon(iconPath);

  if (local) {
    win.webContents.openDevTools();
  }
  win.loadFile(path.join(__dirname, 'dist/octant/index.html'));

  win.webContents.on('did-fail-load', () => {
    win.loadFile(path.join(__dirname, 'dist/octant/index.html'));
  });

  win.webContents.on('new-window', (event, url: string) => {
    event.preventDefault();
    shell.openExternal(url);
  });

  win.on('close', event => {
    const openDialog: boolean = electronStore.get('showDialogue');
    if (openDialog) {
      const messageOptions: MessageBoxOptions = {
        type: 'question',
        buttons: ['Cancel', 'Yes', 'No'],
        defaultId: 2,
        message: 'Do you want to minimize to tray?',
        detail: 'Octant will continue running in the background',
      };

      const result = dialog.showMessageBoxSync(win, messageOptions);
      switch (result) {
        case 0:
          event.preventDefault();
          break;
        case 1:
          event.preventDefault();
          electronStore.set('showDialogue', false);
          win.hide();
          break;
        case 2:
          electronStore.set('showDialogue', false);
          electronStore.set('minimizeToTray', false);
          break;
      }
    } else {
      // @ts-ignore
      const shouldMinimize: boolean = electronStore.get('minimizeToTray');

      if (closing) {
        win = null;
      } else if (shouldMinimize) {
        event.preventDefault();
        win.hide();
      }
    }
  });

  win.on('resize', saveBoundsSoon);
  win.on('move', saveBoundsSoon);

  return win;
}

const startBinary = (port: number) => {
  fs.mkdir(path.join(tmpPath), { recursive: true }, error => {
    if (error) {
      throw error;
    }
  });

  const out = fs.openSync(apiLogPath, 'a');
  const err = fs.openSync(errLogPath, 'a');

  let octantFilename = 'octant';
  if (os.platform() === 'win32') {
    octantFilename = 'octant.exe';
  }

  let serverBinary: string;
  if (local) {
    serverBinary = path.join(__dirname, 'extraResources', octantFilename);
  } else {
    serverBinary = path.join(
      process.resourcesPath,
      'extraResources',
      octantFilename
    );
  }

  const server = child_process.spawn(serverBinary, ['--disable-open-browser'], {
    env: {
      ...process.env,
      NODE_ENV: 'production',
      OCTANT_LISTENER_ADDR: 'localhost:' + port,
    },
    detached: true,
    stdio: ['ignore', out, err],
  });

  serverPid = -server.pid;
  server.unref();
};

try {
  app.on('before-quit', () => {
    process.kill(serverPid, 'SIGHUP');
  });

  app.on('ready', async () => {
    statePath = path.join(app.getPath('userData'), 'state.json');
    loadSavedState();

    const getPort = require('get-port');
    const port = await getPort();
    startBinary(port);
    const w = createWindow();
    w.webContents.on('dom-ready', () => {
      w.webContents.send('port-message', port);
    });

    tray = new TrayMenu(win);

    // In event of a black background issue: https://github.com/electron/electron/issues/15947
    // setTimeout(createWindow, 400);
    session.defaultSession.webRequest.onBeforeSendHeaders(
      { urls: ['ws://localhost:' + port + '/api/v1/stream'] },
      (details, callback) => {
        details.requestHeaders['Origin'] = null;
        callback({ cancel: false, requestHeaders: details.requestHeaders });
      }
    );
  });

  app.on('before-quit', () => {
    closing = true;
  });

  // Quit when all windows are closed.
  app.on('window-all-closed', () => {
    app.quit();
  });

  app.on('activate', () => {
    // On OS X it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
      createWindow();
    } else {
      win.show();
    }
  });
} catch (e) {
  // Catch Error
  // throw e;
}
