/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, OnChanges, SimpleChanges } from '@angular/core';
import { ExtensionView, View } from 'src/app/modules/shared/models/content';
import { SlideInOutAnimation } from './slide-in-out.animation';
import { ResizeEvent } from 'angular-resizable-element';
import { SliderService } from 'src/app/modules/shared/slider/slider.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-slider-view',
  templateUrl: './slider-view.component.html',
  styleUrls: ['./slider-view.component.scss'],
  animations: [SlideInOutAnimation],
})
export class SliderViewComponent
  extends AbstractViewComponent<ExtensionView>
  implements OnChanges {
  style: object = {};
  contentStyle: object = {};
  animationState = 'out';
  contentHeight: number;

  tabs: View[] = [];
  payloads: { [key: string]: string }[] = [];

  constructor(private sliderService: SliderService) {
    super();
  }

  update() {
    const extView = this.v;
    if (extView && extView.config.tabs) {
      this.tabs = [];
      this.payloads = [];
      extView.config.tabs.forEach(tab => {
        this.tabs.push(tab.tab);
        this.payloads.push(tab.payload);
      });
    } else {
      this.animationState = 'out';
      this.sliderService.resetDefault();
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      if (!changes.view.isFirstChange() && this.animationState === 'out') {
        const prevExtView = changes.view.previousValue as ExtensionView;
        if (
          prevExtView.config.tabs === null &&
          changes.view.currentValue.config.tabs.length
        ) {
          this.slide();
          return;
        }
        if (
          changes.view.currentValue.config.tabs.length >
          prevExtView.config.tabs.length
        ) {
          this.slide();
          return;
        }
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
