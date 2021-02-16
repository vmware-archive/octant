// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  EventEmitter,
  Input,
  isDevMode,
  Output,
} from '@angular/core';
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
  readonly defaultItemLimit = 10;
  useSelection = false;
  selectedItem = '';
  url: string;
  position: string;
  action: string;
  itemLimit = this.defaultItemLimit;
  isOpen = false;
  dropdownMenuStyle: object = {};

  @Input() public title: string;

  @Input() public type: string;

  @Input() public items: DropdownItem[];

  @Output() public selectedValue = new EventEmitter<string>();

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

    this.items.forEach(item => {
      if (item.name === view.config.selection) {
        this.selectedItem = item.name;
      }
    });

    this.useSelection = view.config.useSelection;
    if (this.type === 'link') {
      this.url = (view.metadata.title[0] as LinkView).config.ref;
    }

    this.dropdownMenuStyle =
      this.items.length > this.defaultItemLimit ? { 'padding-bottom': 0 } : {};
  }

  identifyItem(index: number, item: DropdownItem): string {
    return item.name;
  }

  toggleShowMore(): void {
    this.itemLimit =
      this.itemLimit === this.items.length
        ? this.defaultItemLimit
        : this.items.length;
  }

  openLink(index): void {
    const item = this.items[index];
    this.selectedItem = item.name;
    this.selectedValue.emit(item.name);
    if (this.useSelection && this.type !== 'icon') {
      this.title = item.label;
    }

    if (this.action) {
      this.websocketService.sendMessage('action.octant.dev/performAction', {
        action: this.action,
        selection: this.selectedItem,
      });
    }

    if (item.url && this.type === 'link') {
      setTimeout(() => {
        this.router.navigateByUrl(item.url);
      }, 0);
    }
    this.isOpen = false;

    if (isDevMode()) {
      console.log('Selected', item.name);
    }
  }
}
