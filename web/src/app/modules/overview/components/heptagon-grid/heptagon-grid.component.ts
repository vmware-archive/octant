// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit } from '@angular/core';
import chunk from 'lodash/chunk';

import { PodStatus } from '../../models/pod-status';
import { Point } from '../../models/point';
import { HoverStatus } from '../heptagon-grid-row/heptagon-grid-row.component';

@Component({
  selector: 'app-heptagon-grid',
  template: `
    <svg [attr.viewBox]="viewBox()">
      <svg:g
        app-heptagon-grid-row
        *ngFor="let row of rows(); let i = index; trackBy: trackByFn"
        class="grid-row"
        [statuses]="row"
        [edgeLength]="edgeLength"
        [row]="i"
        (hoverState)="updateHover($event)"
      />
      <svg:g class="tooltips">
        <svg:g
          app-heptagon-label
          *ngFor="
            let podStatus of podStatuses;
            let i = index;
            trackBy: trackByFn
          "
          class="tooltip"
          [id]="tooltipName(podStatus)"
          [centerPoint]="centerPoint(i)"
          [height]="height()"
          [status]="podStatus"
          [name]="tooltipName(podStatus)"
          [style.opacity]="isActivated(i)"
        />
      </svg:g>
    </svg>
  `,
  styleUrls: ['./heptagon-grid.component.scss'],
})
export class HeptagonGridComponent implements OnInit {
  @Input() podStatuses: PodStatus[] = [];

  @Input()
  edgeLength = 7;

  @Input()
  perRow = 20;

  hoverStates: boolean[][] = [];

  constructor() {}

  ngOnInit() {
    this.rows().forEach(() => {
      this.hoverStates.push([]);
    });
  }

  updateHover(hoverStatus: HoverStatus) {
    this.hoverStates[hoverStatus.row][hoverStatus.col] = hoverStatus.hovered;
  }

  rows() {
    return chunk(this.podStatuses, this.perRow);
  }

  viewBox() {
    const h = this.height();
    return `0 0 ${h * this.perRow * 1.1} ${h * this.rows().length * 1.33 + h}`;
  }

  trackByFn(index, item) {
    return index;
  }

  height() {
    const x = Math.PI / 2 / 7;
    return this.edgeLength / (2 * Math.tan(x));
  }

  centerPoint(index: number) {
    const h = this.height();

    const row = this.row(index);

    let x = h * (index - row * this.perRow) + h / 2;
    const y = h + (row * h + 3);

    const angle = 180 - 90 - 900 / 7 / 2;
    const translateX = (Math.PI / 180) * angle;
    const adjustment = -this.edgeLength * Math.sin(translateX);

    x += adjustment;
    x += h / 2;
    return new Point(x, y);
  }

  tooltipName(status: PodStatus) {
    return `tooltip-${status.name}`;
  }

  row(index: number) {
    return Math.floor(index / this.perRow);
  }

  col(index: number) {
    return index % this.perRow;
  }

  isActivated(index: number) {
    const row = this.row(index);
    const col = this.col(index);

    if (this.hoverStates[row][col]) {
      return 1;
    } else {
      return 0;
    }
  }
}
