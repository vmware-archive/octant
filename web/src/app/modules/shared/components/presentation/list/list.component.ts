// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ChangeDetectionStrategy,
} from '@angular/core';
import {
  LinkView,
  ListView,
  View,
} from 'src/app/modules/shared/models/content';

import { ViewService } from '../../../services/view/view.service';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListComponent implements OnChanges {
  v: ListView;

  @Input() set view(v: View) {
    this.v = v as ListView;
  }
  get view() {
    return this.v;
  }
  title: string;

  iconName: string;

  constructor(private viewService: ViewService) {}

  identifyItem = (index: number, item: View): string => {
    return this.viewService.titleAsText(item.metadata.title);
  };

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      this.title = current.metadata.title
        ? current.metadata.title.map((item: LinkView) => ({
            title: item.config.value,
            url: item.config.ref,
          }))
        : '';

      if (this.v.config.items) {
        this.v.config.items.forEach(item => {
          item.totalItems = this.v.config.items.length;
        });
      }
    }
  }
}
