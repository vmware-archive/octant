// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LinkView } from 'src/app/models/content';

const isUrlAbsolute = url => url.indexOf('://') > 0 || url.indexOf('//') === 0;

@Component({
  selector: 'app-view-link',
  templateUrl: './link.component.html',
  styleUrls: ['./link.component.scss'],
})
export class LinkComponent implements OnChanges {
  @Input() view: LinkView;

  ref: string;
  value: string;
  isAbsolute: boolean;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LinkView;
      this.ref = view.config.ref;
      this.value = view.config.value;
      this.isAbsolute = isUrlAbsolute(this.ref);
    }
  }
}
