// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { LinkView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

const isUrlAbsolute = url =>
  url?.indexOf('://') > 0 || url?.indexOf('//') === 0;

@Component({
  selector: 'app-view-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss'],
})
export class LinkComponent extends AbstractViewComponent<LinkView> {
  ref: string;
  value: string;
  isAbsolute: boolean;
  hasStatus: boolean;

  constructor() {
    super();
  }

  update() {
    const view = this.v;
    this.ref = view.config.ref;
    this.value = view.config.value;
    this.isAbsolute = isUrlAbsolute(this.ref);

    if (view.config.status) {
      this.hasStatus = true;
    }
  }
}
