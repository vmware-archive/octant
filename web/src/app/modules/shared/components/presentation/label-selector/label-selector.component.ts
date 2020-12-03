// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { LabelSelectorView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-label-selector',
  templateUrl: './label-selector.component.html',
  styleUrls: ['./label-selector.component.scss'],
})
export class LabelSelectorComponent extends AbstractViewComponent<LabelSelectorView> {
  key: string;
  value: string;

  constructor() {
    super();
  }

  update() {
    this.key = this.v.config.key;
    this.value = this.v.config.value;
  }
}
