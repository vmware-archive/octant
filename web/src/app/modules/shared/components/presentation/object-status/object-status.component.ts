// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, ViewEncapsulation } from '@angular/core';
import { Node } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-object-status',
  templateUrl: './object-status.component.html',
  styleUrls: ['./object-status.component.scss'],
  encapsulation: ViewEncapsulation.Emulated,
})
export class ObjectStatusComponent {
  @Input() node: Node;

  constructor() {}

  indicatorClass() {
    if (!this.node) {
      return ['progress', 'top', 'success'];
    }

    return [
      'progress',
      'top',
      this.node.status === 'ok' ? 'success' : 'danger',
    ];
  }

  detailsTrackBy(index, item) {
    return index;
  }
}
