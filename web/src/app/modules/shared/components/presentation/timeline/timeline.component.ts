// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { TimelineStep, TimelineView } from '../../../models/content';

@Component({
  selector: 'app-view-timeline',
  templateUrl: './timeline.component.html',
  styleUrls: ['./timeline.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimelineComponent
  extends AbstractViewComponent<TimelineView>
  implements OnInit {
  vertical: boolean;
  steps: TimelineStep[];
  constructor() {
    super();
  }
  update() {
    const view = this.v;
    this.vertical = view.config.vertical;
    this.steps = view.config.steps;
  }
  trackByFn(index, _) {
    return index;
  }
}
