/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app, BrowserWindow, Menu, screen, session, shell } from 'electron';
import * as path from 'path';
import * as child_process from 'child_process';
import * as process from 'process';
import * as os from 'os';
import * as fs from 'fs';

let win: BrowserWindow = null;
let serverPid: any = null;

const args = process.argv.slice(1);
const local = args.some(val => val === '--local');

const template: Electron.MenuItemConstructorOptions[] = [
  {
    label: 'File',
    submenu: [{ label: 'Quit Octant', role: 'close' }],
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
      { role: 'togglefullscreen' },
      { role: 'toggleDevTools' },
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
Menu.setApplicationMenu(menu);

function createWindow(): BrowserWindow {
  const electronScreen = screen;
  const size = electronScreen.getPrimaryDisplay().workAreaSize;

  // Create the browser window.
  win = new BrowserWindow({
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
  });

  win.setIcon(path.join(__dirname, 'dist/dash-frontend/assets/icons/icon.png'));

  if (local) {
    win.webContents.openDevTools();
  }
  win.loadFile(path.join(__dirname, 'dist/dash-frontend/index.html'));

  win.webContents.on('did-fail-load', () => {
    win.loadFile(path.join(__dirname, 'dist/dash-frontend/index.html'));
  });

  // Emitted when the window is closed.
  win.on('closed', () => {
    // Dereference the window object, usually you would store window
    // in an array if your app supports multi windows, this is the time
    // when you should delete the corresponding element.
    win = null;
  });

  return win;
}

const startBinary = (port: number) => {
  const tmpPath = path.join(os.tmpdir(), 'octant');
  fs.mkdir(path.join(tmpPath), { recursive: true }, err => {
    if (err) {
      throw err;
    }
  });

  const out = fs.openSync(path.join(tmpPath, 'api.out.log'), 'a');
  const err = fs.openSync(path.join(tmpPath, 'api.err.log'), 'a');

  let serverBinary: string;
  if (local) {
    serverBinary = path.join(__dirname, 'extraResources', 'octant');
  } else {
    serverBinary = path.join(process.resourcesPath, 'extraResources', 'octant');
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
    const getPort = require('get-port');
    const port = await getPort();
    startBinary(port);
    const w = createWindow();
    w.webContents.on('dom-ready', () => {
      w.webContents.send('port-message', port);
    });

    // In event of a black background issue: https://github.com/electron/electron/issues/15947
    //setTimeout(createWindow, 400);
    session.defaultSession.webRequest.onBeforeSendHeaders(
      { urls: ['ws://localhost:' + port + '/api/v1/stream'] },
      (details, callback) => {
        details.requestHeaders['Origin'] = null;
        callback({ cancel: false, requestHeaders: details.requestHeaders });
      }
    );
  });

  // Quit when all windows are closed.
  app.on('window-all-closed', () => {
    // On OS X it is common for applications and their menu bar
    // to stay active until the user quits explicitly with Cmd + Q
    if (process.platform !== 'darwin') {
      app.quit();
    }
  });

  app.on('activate', () => {
    // On OS X it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (win === null) {
      createWindow();
    }
  });
} catch (e) {
  // Catch Error
  // throw e;
}
