// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { PodStatusView } from '../../../models/content';
import { PodStatus } from '../../../models/pod-status';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-pod-status',
  templateUrl: './pod-status.component.html',
  styleUrls: ['./pod-status.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class PodStatusComponent extends AbstractViewComponent<PodStatusView> {
  @ViewChild('container') private container: ElementRef;

  edgeLength = 7;

  podStatuses: PodStatus[] = [];

  constructor() {
    super();
  }

  update() {
    const pods = this.v.config.pods;

    if (pods) {
      this.podStatuses = Object.keys(pods)
        .sort()
        .map((podName: string): PodStatus => {
          return {
            name: podName,
            status: pods[podName].status,
          };
        });
    }
  }
}
