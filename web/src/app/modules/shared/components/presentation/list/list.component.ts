// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ListView, View } from 'src/app/modules/shared/models/content';

import { IconService } from '../../../services/icon/icon.service';
import { ViewService } from '../../../services/view/view.service';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
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

  constructor(
    private iconService: IconService,
    private viewService: ViewService
  ) {}

  identifyItem = (index: number, item: View): string => {
    return this.viewService.titleAsText(item.metadata.title);
  };

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      this.title = this.viewService.viewTitleAsText(current);
      this.iconName = this.iconService.load(current.config);
    }
  }
}
