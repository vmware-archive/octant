// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';

import { PodStatus } from '../../../../../modules/overview/models/pod-status';
import { Point } from '../../../../../modules/overview/models/point';

export interface HoverStatus {
  row: number;
  col: number;
  hovered: boolean;
}

@Component({
  selector: '[app-heptagon-grid-row]',
  template: `
    <svg:g
      *ngFor="let status of statuses; let i = index; trackBy: trackByFn"
      app-heptagon
      [id]="name(status)"
      class="row"
      [status]="status"
      [centerPoint]="centerPoint(i)"
      [edgeLength]="edgeLength"
      [isFlipped]="isFlipped(i)"
      (hovered)="updateHover($event, i)"
    />
  `,
  styleUrls: ['./heptagon-grid-row.component.scss'],
})
export class HeptagonGridRowComponent implements OnInit {
  @Input()
  statuses: PodStatus[];

  @Input()
  edgeLength: number;

  @Input()
  row: number;

  @Output()
  hoverState = new EventEmitter<HoverStatus>();

  constructor() {}

  ngOnInit() {}

  updateHover(hovered, index) {
    this.hoverState.emit({
      row: this.row,
      col: index,
      hovered,
    });
  }

  centerPoint(index: number) {
    const h = this.height();

    let x = h * index + h / 2;
    const y = h + (this.row * h + 3);

    const angle = 180 - 90 - 900 / 7 / 2;
    const translateX = (Math.PI / 180) * angle;
    const adjustment = -this.edgeLength * Math.sin(translateX);

    x += adjustment;
    x += h / 2;
    return new Point(x, y);
  }

  isFlipped(index: number) {
    return index % 2 !== 0;
  }

  trackByFn(index, item) {
    return index;
  }

  height() {
    const x = Math.PI / 2 / 7;
    return this.edgeLength / (2 * Math.tan(x));
  }

  name(status: PodStatus) {
    return `heptagon-${status.name}`;
  }
}
