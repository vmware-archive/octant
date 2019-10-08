// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ExpressionSelectorView } from 'src/app/models/content';

@Component({
  selector: 'app-view-expression-selector',
  templateUrl: './expression-selector.component.html',
  styleUrls: ['./expression-selector.component.scss'],
})
export class ExpressionSelectorComponent implements OnChanges {
  @Input() view: ExpressionSelectorView;

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
