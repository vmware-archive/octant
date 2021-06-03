// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, EventEmitter, Output } from '@angular/core';
import { TabsView, View } from '../../../models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

interface Tab {
  name: string;
  view: View;
  accessor: string;
}

@Component({
  selector: 'app-tabs-view',
  templateUrl: './tabs-view.component.html',
  styleUrls: ['./tabs-view.component.scss'],
})
export class TabsViewComponent extends AbstractViewComponent<TabsView> {
  activeTab: string;
  tabs: View[] = [];
  orientation: string;

  constructor() {
    super();
  }

  update() {
    this.tabs = this.v.config.tabs;
    this.orientation = this.v.config?.orientation || 'horizontal';
  }

  clickTab(tabAccessor: string) {
    if (this.activeTab === tabAccessor) {
      return;
    }
    this.activeTab = tabAccessor;
  }

  identifyTab(index: number, item: Tab): string {
    return item.name;
  }
}
