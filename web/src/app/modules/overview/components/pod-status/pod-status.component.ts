// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';
import { PodStatusView } from '../../../../models/content';
import { PodStatus } from '../../models/pod-status';

@Component({
  selector: 'app-pod-status',
  templateUrl: './pod-status.component.html',
  styleUrls: ['./pod-status.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class PodStatusComponent implements OnChanges {
  @ViewChild('container', { static: false }) private container: ElementRef;

  @Input() view: PodStatusView;

  edgeLength = 7;

  podStatuses: PodStatus[] = [];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    const pods = changes.view.currentValue.config.pods;

    const statuses = Object.keys(pods)
      .sort()
      .map(
        (podName: string): PodStatus => {
          return {
            name: podName,
            status: pods[podName].status,
          };
        }
      );

    this.podStatuses = statuses;
  }
}
