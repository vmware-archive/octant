// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { ExpressionSelectorView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-expression-selector',
  templateUrl: './expression-selector.component.html',
  styleUrls: ['./expression-selector.component.scss'],
})
export class ExpressionSelectorComponent extends AbstractViewComponent<ExpressionSelectorView> {
  key: string;
  operator: string;
  values: string;

  constructor() {
    super();
  }

  update() {
    const view = this.v;
    this.key = view.config.key;
    this.operator = view.config.operator;
    this.values = view.config.values?.join('|');
  }
}
