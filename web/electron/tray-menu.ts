/*
 *  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app, BrowserWindow, Menu, Tray, nativeImage, shell, MenuItem } from 'electron';
import { errLogPath, iconPath, greyIconPath } from './paths';
import * as WebSocket from 'ws';
import * as open from 'open';

export class TrayMenu {
  public readonly tray: Tray;
  public wsConn : WebSocket;
  private menuState : {
    contexts: string[],
    namespaces: string[],
    currentContext: (string | null),
    buildInfo: {version: string, commit: string, time: string}
  }

  constructor(public window: BrowserWindow, public websocketUrl: string) {
    this.tray = new Tray(this.createNativeImage());
    this.wsConn = new WebSocket(this.websocketUrl);
    this.menuState = { contexts: [], namespaces: [], currentContext: null, 
      buildInfo: {version: '', commit: '', time: ''}
    }

    this.setMenu();
    this.startOctantEventListener();
  }
  
  createNativeImage(): Electron.NativeImage {
    const image = nativeImage.createFromPath(greyIconPath);
    image.setTemplateImage(true);
    return image.resize({ width: 16, height: 16 });
  }

  startOctantEventListener() {
    this.wsConn.on('message', (msg) => {
      const config = JSON.parse(`${msg}`);

      switch (config.type) {
        case 'event.octant.dev/kubeConfig':
          this.menuState.contexts = config.data.contexts.map(c => c.name);
          this.menuState.currentContext = config.data.currentContext;
          this.setMenu();
          break;
        case 'event.octant.dev/buildInfo':
          this.menuState.buildInfo = config.data;
          this.setAboutOptions();
          this.setMenu();
          break;
      }
    });
  }

  setMenu() {
    let menu = new Menu();
    menu.append(new MenuItem({label: 'Open Octant', type: 'normal', click: () => this.window.show()}));
    menu.append(new MenuItem({label: 'View Logs', type: 'normal', click: () => shell.showItemInFolder(errLogPath)}));
    
    menu.append(this.contextSubMenu());
    menu.append(this.aboutOctantSubMenu());
    menu.append(this.openIssueMenuItem());
    menu.append(this.octantDocsMenuItem());
    
    menu.append(new MenuItem({label: 'Quit', type: 'normal', click: () => app.quit()}))
    this.tray.setContextMenu(menu);
  }

  contextSubMenu() : MenuItem {
    let contextMenu = new Menu();
    const cItems = this.menuState.contexts.forEach((context) => {
      contextMenu.append(new MenuItem({
        label: context, 
        type: 'checkbox', 
        checked: context === this.menuState.currentContext,
        click: () => { this.wsConn.send(JSON.stringify({
          type: 'action.octant.dev/setContext',
          payload: { requestedContext: context }
        })) }
      }));
    });
    

    return new MenuItem({label: 'Contexts', type: 'submenu', submenu: contextMenu});
  }

  aboutOctantSubMenu() : MenuItem {
    return new MenuItem({label: 'About Octant', role: 'about'});
  }

  openIssueMenuItem() : MenuItem {
    const newIssueLink = 'https://github.com/vmware-tanzu/octant/issues/new/choose';
    return new MenuItem({label: 'Open an Issue/Provide Feedback', type: 'normal', click: () => open(newIssueLink)});
  }

  octantDocsMenuItem() : MenuItem {
    const docsLink = 'https://octant.dev/';
    return new MenuItem({label: 'Octant Documentation', type: 'normal', click: () => open(docsLink)});
  }

  setAboutOptions() : void {
    app.setAboutPanelOptions({
      applicationName: 'Octant',
      applicationVersion: `Version: ${this.menuState.buildInfo.version}\nCommit: ${this.menuState.buildInfo.commit}\nBuilt: ${this.menuState.buildInfo.time}`,
      website: 'https://octant.dev',
      iconPath: iconPath
    })
  }
}
