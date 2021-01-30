/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app, BrowserWindow, Menu, screen, session } from 'electron';
import { ApplicationMenu } from './electron/application-menu';
import { TrayMenu } from './electron/tray-menu';
import { apiLogPath, errLogPath, tmpPath, iconPath } from './electron/paths';
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

<<<<<<< HEAD
const tmpPath = path.join(os.tmpdir(), 'octant');
const apiLogPath = path.join(tmpPath, 'api.out-' + date + '.log');
const errLogPath = path.join(tmpPath, 'api.err-' + date + '.log');

const template: Electron.MenuItemConstructorOptions[] = [
  {
    label: 'File',
    submenu: [
      { label: 'Preferences',
        click(){
          win.webContents.send('openPreferences');
        } },
      { type: 'separator' },
      { label: 'Quit Octant', role: 'close' }
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
      { role: 'togglefullscreen' },
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
Menu.setApplicationMenu(menu);
=======
const applicationMenu = new ApplicationMenu();
Menu.setApplicationMenu(applicationMenu.menu);
>>>>>>> Add system menu tray and minimize/hide logic

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

  var options = {
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

  win.on('close', (event) => {
    if (closing) {
      win = null;
    } else {
      event.preventDefault();
      win.hide();
    }
  });

  // Emitted when the window is closed.
  // win.on('closed', () => {
  //   // Dereference the window object, usually you would store window
  //   // in an array if your app supports multi windows, this is the time
  //   // when you should delete the corresponding element.
  //   win = null;
  // });
  win.on('resize', saveBoundsSoon);
  win.on('move', saveBoundsSoon);

  return win;
}

const startBinary = (port: number) => {
  fs.mkdir(path.join(tmpPath), { recursive: true }, err => {
    if (err) {
      throw err;
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
    //setTimeout(createWindow, 400);
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
  })

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
