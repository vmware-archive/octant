// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  Input,
  ViewChild,
  AfterViewInit,
  ChangeDetectorRef,
  AfterViewChecked,
  ChangeDetectionStrategy,
} from '@angular/core';
import { PodStatus } from '../../models/pod-status';
import { Point } from '../../models/point';

@Component({
  selector: '[app-heptagon-label]',
  template: `
    <svg:rect
      #container
      [attr.y]="y()"
      [attr.x]="containerX()"
      [attr.width]="containerWidth()"
      [attr.height]="containerHeight()"
    />
    <svg:text
      #label
      [attr.y]="y()"
      [attr.x]="containerX()"
      [attr.dx]="textPaddingX()"
      [attr.dy]="textPaddingY()"
      text-anchor="left"
    >
      {{ status.name }}
    </svg:text>
  `,
  styleUrls: ['./heptagon-label.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HeptagonLabelComponent implements AfterViewChecked {
  @ViewChild('container')
  container: ElementRef;

  @ViewChild('label')
  labelText: ElementRef;

  @Input()
  centerPoint: Point;

  @Input()
  height: number;

  @Input()
  status: PodStatus;

  @Input()
  padding = {
    x: 15,
    y: 10,
  };

  @Input()
  name: string;

  constructor(private cd: ChangeDetectorRef) {}

  ngAfterViewChecked(): void {
    this.cd.detectChanges();
  }

  x() {
    return this.centerPoint.x - this.height / 2;
  }

  y() {
    return this.centerPoint.y - this.height / 2 - 5;
  }

  containerX() {
    return this.x() - this.width() - 5;
  }

  containerWidth() {
    return `${this.width()}px`;
  }

  containerHeight() {
    return this.bBox().height + 2 * this.padding.y;
  }

  width() {
    return this.bBox().width + 2 * this.padding.x;
  }

  bBox() {
    return this.labelText.nativeElement.getBBox();
  }

  fontSize() {
    const el = this.labelText.nativeElement;
    const style = window.getComputedStyle(el, null);
    const fontSizeRaw = style.getPropertyValue('font-size');
    return parseFloat(fontSizeRaw);
  }

  textPaddingX() {
    return `${this.padding.x}px`;
  }

  textPaddingY() {
    return `${this.bBox().height}px`;
  }
}
