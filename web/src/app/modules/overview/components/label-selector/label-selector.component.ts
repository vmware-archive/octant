// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LabelSelectorView } from 'src/app/models/content';

@Component({
  selector: 'app-view-label-selector',
  templateUrl: './label-selector.component.html',
  styleUrls: ['./label-selector.component.scss'],
})
export class LabelSelectorComponent implements OnChanges {
  @Input() view: LabelSelectorView;

  key: string;
  value: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LabelSelectorView;
      this.key = view.config.key;
      this.value = view.config.value;
    }
  }
}
