/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, Input } from '@angular/core';
import { DonutChartView, DonutSegment } from '../../../models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { Point } from '../../../models/point';

export interface SegmentDescriptor {
  id: string;
  fillColor: string;
  path: string;
}

@Component({
  selector: 'app-view-donut-chart',
  templateUrl: './donut-chart.component.html',
  styleUrls: ['./donut-chart.component.scss'],
})
export class DonutChartComponent extends AbstractViewComponent<DonutChartView> {
  @Input() circumference = 100;
  scale: string;

  segments: SegmentDescriptor[] = [];
  tooltipText = '';
  tooltipX = '';
  tooltipY = '';
  selectedSegment = -1;

  donutPadding = 8;
  svgSize = 250;
  donutThickness = 0.2;

  constructor() {
    super();
  }

  update() {
    if (this.v?.config?.size) {
      this.scale = String(this.v.config.size) + '%';
      if (this.v.config.thickness) {
        this.donutThickness = this?.v?.config?.thickness / 100;
      }
      this.createSegments();
    }
  }

  trackByDescriptor(index: number, item: SegmentDescriptor) {
    return item?.id;
  }

  viewBox(): string {
    return `0 0 ${this.svgSize} ${this.svgSize}`;
  }

  itemCount(): number {
    if (!this.v) {
      return 0;
    }
    return this?.v?.config?.segments?.reduce<number>(
      (a: number, s: DonutSegment) => a + s.count,
      0
    );
  }

  itemLabel(): string {
    if (!this.v) {
      return '';
    }

    return this.itemCount() > 1
      ? this?.v?.config?.labels?.plural
      : this?.v?.config?.labels?.singular;
  }

  createSegments(): void {
    const donutRadius = (this.svgSize - this.donutPadding) / 2;
    let startAngle = 0;
    this.segments = [];

    this.v.config.segments
      .sort((a, b) =>
        a.status || b.status
          ? a.status > b.status
            ? 1
            : -1
          : b.count > a.count
          ? 1
          : -1
      )
      .forEach((segment, index) => {
        const segmentTotal = this.itemCount();
        const center = this.svgSize / 2;
        const cutoutRadius = segment.thickness
          ? donutRadius * (1 - segment.thickness / 100)
          : this.hasCriticalSegments()
          ? Math.max(0, donutRadius * (1 - (3 * this.donutThickness) / 2))
          : donutRadius * (1 - this.donutThickness);
        const startRadius =
          this.hasCriticalSegments() && !this.isCriticalSegment(segment)
            ? Math.max(0, donutRadius * (1 - (5 * this.donutThickness) / 4))
            : Math.max(0, cutoutRadius);
        const endRadius =
          this.hasCriticalSegments() && !this.isCriticalSegment(segment)
            ? donutRadius * (1 - this.donutThickness / 4)
            : donutRadius;
        const segmentAngle = (2 * Math.PI * segment.count) / segmentTotal;
        const finalAngle = startAngle + segmentAngle - 0.001;
        const largeArc =
          (finalAngle - startAngle) % (Math.PI * 2) > Math.PI ? 1 : 0;
        const { p1: start1, p2: end1 } = this.getArcPoints(
          center,
          startAngle,
          startRadius,
          endRadius
        );
        const { p1: start2, p2: end2 } = this.getArcPoints(
          center,
          finalAngle,
          startRadius,
          endRadius
        );
        const path = `M ${start1.x} ${start1.y} A ${endRadius} ${endRadius} 0 ${largeArc} 1 ${start2.x} ${start2.y} L ${end2.x} ${end2.y} A ${startRadius}, ${startRadius} 0 ${largeArc} 0 ${end1.x} ${end1.y} Z`;
        this.segments.push({
          path,
          id: `path${index}`,
          fillColor: this.getSegmentColor(segment),
        });
        startAngle += segmentAngle;
      });
  }

  getArcPoints(center, angle, radius1, radius2): { p1: Point; p2: Point } {
    const centerPoint = new Point(center, center);

    const p1 = centerPoint.projectRadian({ magnitude: radius2, angle });
    const p2 = centerPoint.projectRadian({ magnitude: radius1, angle });
    return { p1, p2 };
  }

  mouseEnterEvent(event: any, segmentIndex: number) {
    if (segmentIndex !== this.selectedSegment) {
      this.selectedSegment = segmentIndex;
      const segment = this.v?.config?.segments[this.selectedSegment];
      if (segment) {
        const scrollOffset = this.getScrollOffset(event.target.ownerSVGElement);
        const label =
          segment.count > 1
            ? this?.v?.config?.labels?.plural
            : this?.v?.config?.labels?.singular;

        this.tooltipX = `${event.offsetX}px`;
        this.tooltipY = `${
          event.offsetY -
          event.target.ownerSVGElement.clientHeight -
          scrollOffset -
          48
        }px`;
        this.tooltipText = segment.description
          ? segment.description
          : `${segment.count} ${label} with status ${segment.status}`;
      }
    }
  }

  mouseLeaveEvent() {
    this.selectedSegment = -1;
  }

  getScrollOffset(element: any) {
    if (!element) return 0;
    return element.scrollTop + this.getScrollOffset(element.parentElement);
  }

  hasCriticalSegments(): boolean {
    let hasSegments = false;
    this?.v?.config?.segments?.forEach(segment => {
      if (this.isCriticalSegment(segment)) {
        hasSegments = true;
      }
    });
    return hasSegments;
  }

  isCriticalSegment(segment: DonutSegment): boolean {
    return segment.status === 'error' || segment.status === 'warning';
  }

  getSegmentColor(segment: DonutSegment): string {
    if (segment.color) {
      return segment.color;
    }
    switch (segment.status) {
      case 'ok':
        return '#60b515';
      case 'warning':
        return '#f57600';
      case 'error':
        return '#e12200';
      default:
        return 'd2d3d4';
    }
  }
}
