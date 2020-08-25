// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import { LoadingView, View } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-loading',
  templateUrl: './loading.component.html',
  styleUrls: ['./loading.component.scss'],
})
export class LoadingComponent extends AbstractViewComponent<LoadingView> {
  value: string;

  constructor() {
    super();
  }

  update() {
    this.value = this.v.config.value;
  }
}
