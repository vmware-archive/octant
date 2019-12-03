/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component, Input } from '@angular/core';
import { View, ExtensionView } from 'src/app/models/content';
import { SlideInOutAnimation } from './slide-in-out.animation';

@Component({
  selector: 'app-slider-view',
  templateUrl: './slider-view.component.html',
  styleUrls: ['./slider-view.component.scss'],
  animations: [SlideInOutAnimation],
})
export class SliderViewComponent {
  @Input() view: ExtensionView;

  animationState = 'out';

  slide() {
    this.animationState = this.animationState === 'out' ? 'in' : 'out';
  }
}
