// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ListView, View } from 'src/app/models/content';

import { IconService } from '../../services/icon.service';
import { ViewService } from '../../services/view/view.service';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ListComponent implements OnChanges {
  @Input() listView: ListView;
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
    if (changes.listView) {
      const current = changes.listView.currentValue;
      this.title = this.viewService.viewTitleAsText(current);
      this.iconName = this.iconService.load(current.config);
    }
  }
}
