// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
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

  items: View[];

  private previous: string;

  constructor(
    private iconService: IconService,
    private viewService: ViewService,
    private cdr: ChangeDetectorRef
  ) {
    super();
  }

  identifyItem = (index: number, item: View): string => {
    return this.viewService.titleAsText(item.metadata.title);
  };

  update() {
    const current = this.v;

    const cur = JSON.stringify(current);
    if (current.config.items && cur !== this.previous) {
      this.items = this.v.config.items;
      this.initialChildCount = this.items.length;
      this.items.forEach(item => {
        item.totalItems = current.config.items.length;
      });
      this.previous = cur;
      this.cdr.markForCheck();
    }
  }
}
