// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LabelSelectorView, View } from 'src/app/shared/models/content';

@Component({
  selector: 'app-view-label-selector',
  templateUrl: './label-selector.component.html',
  styleUrls: ['./label-selector.component.scss'],
})
export class LabelSelectorComponent implements OnChanges {
  private v: LabelSelectorView;

  @Input() set view(v: View) {
    this.v = v as LabelSelectorView;
  }
  get view() {
    return this.v;
  }
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
