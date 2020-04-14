/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ExtensionView, View } from '../../../../shared/models/content';
import { SlideInOutAnimation } from './slide-in-out.animation';
import { ResizeEvent } from 'angular-resizable-element';
import { SliderService } from '../../../../shared/slider/slider.service';

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

  constructor(private sliderService: SliderService) {}

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

        if (!changes.view.isFirstChange() && this.animationState === 'out') {
          const prevExtView = changes.view.previousValue as ExtensionView;
          if (prevExtView.config.tabs === null && extView.config.tabs.length) {
            this.slide();
            return;
          }
          if (extView.config.tabs.length > prevExtView.config.tabs.length) {
            this.slide();
            return;
          }
        }
      } else {
        this.animationState = 'out';
        this.sliderService.resetDefault();
      }
    }
  }

  slide() {
    this.animationState = this.animationState === 'out' ? 'in' : 'out';

    // Note: these checks represent state after click
    if (this.animationState === 'in' && this.contentHeight) {
      Object.assign(this.style, { height: `${this.contentHeight}px` });
      this.sliderService.setHeight(this.contentHeight);
      return;
    }
    if (this.animationState === 'in' && !this.contentHeight) {
      this.sliderService.setHeight(288);
      return;
    }

    if (this.animationState === 'out') {
      this.sliderService.resetDefault();
      this.style = {};
      return;
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
      this.sliderService.setHeight(this.contentHeight);
    }
  }

  onResize(event) {
    Object.assign(this.style, { top: `inherit` });
  }
}
