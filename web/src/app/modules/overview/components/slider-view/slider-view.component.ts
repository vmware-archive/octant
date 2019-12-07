/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, Input } from '@angular/core';
import { ExtensionView } from 'src/app/models/content';
import { SlideInOutAnimation } from './slide-in-out.animation';
import { ResizeEvent } from 'angular-resizable-element';

@Component({
  selector: 'app-slider-view',
  templateUrl: './slider-view.component.html',
  styleUrls: ['./slider-view.component.scss'],
  animations: [SlideInOutAnimation],
})
export class SliderViewComponent {
  @Input() view: ExtensionView;

  style: object = {};
  animationState = 'out';

  slide() {
    this.animationState = this.animationState === 'out' ? 'in' : 'out';
    this.style = {};
  }

  onResizeTop(event: ResizeEvent): void {
    if (this.animationState === 'in') {
      this.style = {
        top: `${event.rectangle.top}px`,
        height: `${event.rectangle.height}px`,
        cursor: `ns-resize`,
      };

      if (this.style !== {}) {
        this.animationState = 'in';
      }
    }
  }
}
