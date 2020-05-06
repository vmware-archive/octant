// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LinkView, View } from 'src/app/modules/shared/models/content';

const isUrlAbsolute = url => url.indexOf('://') > 0 || url.indexOf('//') === 0;

@Component({
  selector: 'app-view-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss'],
})
export class LinkComponent implements OnChanges {
  private v: LinkView;

  @Input() set view(v: View) {
    this.v = v as LinkView;
  }
  get view() {
    return this.v;
  }

  ref: string;
  value: string;
  isAbsolute: boolean;
  hasStatus: boolean;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LinkView;
      this.ref = view.config.ref;
      this.value = view.config.value;
      this.isAbsolute = isUrlAbsolute(this.ref);

      if (view.config.status) {
        this.hasStatus = true;
      }
    }
  }
}
