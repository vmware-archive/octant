// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { ChangeDetectionStrategy, Component } from '@angular/core';
import {
  LinkView,
  ListView,
  View,
} from 'src/app/modules/shared/models/content';

import { IconService } from '../../../services/icon/icon.service';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListComponent extends AbstractViewComponent<ListView> {
  title: { title: string; url: string }[];

  iconName: string;

  constructor(
    private iconService: IconService,
    private viewService: ViewService
  ) {
    super();
  }

  identifyItem = (index: number, item: View): string => {
    return this.viewService.titleAsText(item.metadata.title);
  };

  update() {
    const current = this.v;
    this.title = current.metadata.title
      ? current.metadata.title.map((item: LinkView) => ({
          title: item.config.value,
          url: item.config.ref,
        }))
      : [];

    if (this.v.config.items) {
      this.initialChildCount = this.v.config.items.length;
      this.v.config.items.forEach(item => {
        item.totalItems = this.v.config.items.length;
      });
    }
  }
}
