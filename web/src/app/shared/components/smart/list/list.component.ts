// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ListView, View } from 'src/app/shared/models/content';

import { IconService } from '../../../../modules/overview/services/icon.service';

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
  title: View[];

  iconName: string;

  constructor(private iconService: IconService) {}

  identifyItem = (index: number, _: View): number => {
    return index;
  };

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      this.title = current.metadata.title;
      this.iconName = this.iconService.load(current.config);
    }
  }
}
