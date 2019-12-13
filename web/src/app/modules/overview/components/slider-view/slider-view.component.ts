/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, Input } from '@angular/core';
import { ExtensionView, View } from 'src/app/models/content';
import { SlideInOutAnimation } from './slide-in-out.animation';
import { ResizeEvent } from 'angular-resizable-element';
import { OnChanges, SimpleChanges } from '@angular/core';

@Component({
  selector: 'app-slider-view',
  templateUrl: './slider-view.component.html',
  styleUrls: ['./slider-view.component.scss'],
  animations: [SlideInOutAnimation],
})
export class SliderViewComponent implements OnChanges {
  @Input() view: ExtensionView;

  style: object = {};
  contentStyle: object = {};
  animationState = 'out';
  contentHeight: number;

  tabs: View[] = [];
  payloads: { [key: string]: string }[] = [];

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const extView = changes.view.currentValue as ExtensionView;
      if (extView.config.tabs) {
        this.tabs = [];
        this.payloads = [];
        extView.config.tabs.forEach(tab => {
          this.tabs.push(tab.tab);
          this.payloads.push(tab.payload);
        });
      } else {
        this.animationState = 'out';
      }
    }
  }

  slide() {
    this.animationState = this.animationState === 'out' ? 'in' : 'out';
    this.style = {};

    if (this.contentHeight) {
      Object.assign(this.style, { height: `${this.contentHeight}px` });
    }
  }

  onResizeTop(event: ResizeEvent): void {
    if (this.animationState === 'in') {
      this.style = {
        top: `${event.rectangle.top}px`,
        height: `${event.rectangle.height}px`,
        cursor: `ns-resize`,
      };

      this.contentStyle = {
        height: `${event.rectangle.height}px`,
      };
      this.contentHeight = event.rectangle.height;
    }
  }
}
