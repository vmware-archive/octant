// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, isDevMode } from '@angular/core';
import {
  DropdownItem,
  DropdownView,
  LinkView,
} from 'src/app/modules/shared/models/content';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { Router } from '@angular/router';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';

@Component({
  selector: 'app-view-dropdown',
  templateUrl: './dropdown.component.html',
  styleUrls: ['./dropdown.component.scss'],
})
export class DropdownComponent extends AbstractViewComponent<DropdownView> {
  useSelection = false;
  selectedItem = '';
  title: string;
  url: string;
  position: string;
  type: string;
  action: string;
  items: DropdownItem[];

  constructor(
    private viewService: ViewService,
    private websocketService: WebsocketService,
    private router: Router
  ) {
    super();
  }

  update() {
    const view = this.v;
    this.title = this.viewService.viewTitleAsText(view);
    this.position = view.config.position;
    this.type = view.config.type;
    this.action = view.config.action;
    this.items = view.config.items;

    this.useSelection = view.config.useSelection;
    if (this.type === 'link') {
      this.url = (view.metadata.title[0] as LinkView).config.ref;
    }
  }

  identifyItem(index: number, item: DropdownItem): string {
    return item.name;
  }

  openLink(index): void {
    const item = this.items[index];
    this.selectedItem = item.name;
    if (this.useSelection && this.type !== 'icon') {
      this.title = item.label;
    }

    this.websocketService.sendMessage('action.octant.dev/performAction', {
      action: this.action,
      selection: this.selectedItem,
    });

    if (item.url) {
      this.router.navigateByUrl(item.url);
    }

    if (isDevMode()) {
      console.log('Selected', item.name);
    }
  }
}
