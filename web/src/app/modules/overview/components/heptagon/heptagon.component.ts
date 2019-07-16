// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterViewInit,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnInit,
  Output,
} from '@angular/core';

import { PodStatus } from '../../models/pod-status';
import { Point } from '../../models/point';
import { Vector } from '../../models/vector';
import { Router } from '@angular/router';

@Component({
  selector: '[heptagon]',
  template: `
    <svg:path [attr.d]="path()" [ngClass]="style()" (click)="navigate()" />
  `,
  styleUrls: ['./heptagon.component.scss'],
})
export class HeptagonComponent implements OnInit, AfterViewInit {
  @Input()
  status: PodStatus;

  @Input()
  centerPoint: Point;

  @Input()
  edgeLength: number;

  @Input()
  isFlipped = false;

  @Output()
  hovered = new EventEmitter<boolean>();

  constructor(private elementRef: ElementRef, private router: Router) {}

  ngOnInit() {}

  ngAfterViewInit(): void {
    const el = this.elementRef.nativeElement;

    el.addEventListener('mouseover', () => {
      this.hovered.emit(true);
    });

    el.addEventListener('mouseout', () => {
      this.hovered.emit(false);
    });
  }

  label() {
    return this.status.name;
  }

  style() {
    return `heptagon status-${this.status.status}`;
  }

  path() {
    return this.points()
      .map((point, index) => {
        const command = index === 0 ? 'M' : 'L';
        return `${command}${point.toString()}`;
      })
      .join(' ');
  }

  points() {
    return Array(7)
      .fill({})
      .map((_, index) => {
        const radian = ((Math.PI / 180) * 360) / 14;
        const vector: Vector = {
          angle: (360 / 7) * index + 360 / 28,
          magnitude: (this.edgeLength * 0.5) / Math.sin(radian),
        };

        if (this.isFlipped) {
          vector.angle += 360 / 14;
        }

        const projected = this.centerPoint.project(vector);
        if (this.isFlipped) {
          const translateY1 = ((Math.PI / 180) * 360) / 7;
          let adjustment = -this.edgeLength * Math.sin(translateY1);

          const translateY2 = ((Math.PI / 180) * 360) / 14;
          adjustment = -this.edgeLength * Math.sin(translateY2);

          projected.y += adjustment;
        }

        return projected;
      });
  }

  navigate() {
    this.router.navigate([
      '/content',
      'overview',
      'namespace',
      'default',
      'workloads',
      'pods',
      this.status.name,
    ]);
  }
}
