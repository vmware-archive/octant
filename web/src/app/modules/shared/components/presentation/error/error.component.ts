// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { ErrorView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-error',
  templateUrl: './error.component.html',
  styleUrls: ['./error.component.scss'],
})
export class ErrorComponent extends AbstractViewComponent<ErrorView> {
  source: string;

  constructor() {
    super();
  }

  update() {
    this.source = this.v.config.data;
  }
}
