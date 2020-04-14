// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import {
  ExpressionSelectorView,
  LabelSelectorView,
  SelectorsView,
  View,
} from '../../../../shared/models/content';

@Component({
  selector: 'app-view-selectors',
  templateUrl: './selectors.component.html',
  styleUrls: ['./selectors.component.scss'],
})
export class SelectorsComponent {
  v: SelectorsView;

  @Input() set view(v: View) {
    this.v = v as SelectorsView;
  }
  get view() {
    return this.v;
  }

  identifyItem(
    index: number,
    item: ExpressionSelectorView | LabelSelectorView
  ): string {
    const { key } = item.config;
    const labelSelector = item as LabelSelectorView;
    const expressionSelector = item as ExpressionSelectorView;
    if (labelSelector.config.value) {
      return `${key}-${labelSelector.config.value}`;
    } else if (expressionSelector.config.values) {
      return `${key}-${
        expressionSelector.config.operator
      }-${expressionSelector.config.values.join(',')}`;
    }
  }
}
