// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import {
  ExpressionSelectorView,
  LabelSelectorView,
  SelectorsView,
} from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-selectors',
  templateUrl: './selectors.component.html',
  styleUrls: ['./selectors.component.scss'],
})
export class SelectorsComponent extends AbstractViewComponent<SelectorsView> {
  constructor() {
    super();
  }

  update() {}

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
