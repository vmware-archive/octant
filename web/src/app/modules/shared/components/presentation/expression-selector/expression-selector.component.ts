// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import {
  ExpressionSelectorView,
  View,
} from '../../../../shared/models/content';

@Component({
  selector: 'app-view-expression-selector',
  templateUrl: './expression-selector.component.html',
  styleUrls: ['./expression-selector.component.scss'],
})
export class ExpressionSelectorComponent implements OnChanges {
  private v: ExpressionSelectorView;

  @Input() set view(v: View) {
    this.v = v as ExpressionSelectorView;
  }
  get view() {
    return this.v;
  }

  key: string;
  operator: string;
  values: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as ExpressionSelectorView;
      this.key = view.config.key;
      this.operator = view.config.operator;
      this.values = view.config.values.join(',');
    }
  }
}
