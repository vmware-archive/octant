/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component } from '@angular/core';
import { SlideInOutAnimation } from './slide-in-out.animation';

@Component({
  selector: 'app-slider-view',
  templateUrl: './slider-view.component.html',
  styleUrls: ['./slider-view.component.scss'],
  animations: [SlideInOutAnimation],
})
export class SliderViewComponent {
  animationState = 'out';

  slide() {
    this.animationState = this.animationState === 'out' ? 'in' : 'out';
  }
}
