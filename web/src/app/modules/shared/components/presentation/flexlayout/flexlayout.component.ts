// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import {
  ButtonGroupView,
  FlexLayoutItem,
  FlexLayoutView,
  View,
} from '../../../../shared/models/content';
import trackByIndex from '../../../../../util/trackBy/trackByIndex';

@Component({
  selector: 'app-view-flexlayout',
  templateUrl: './flexlayout.component.html',
  styleUrls: ['./flexlayout.component.scss'],
})
export class FlexlayoutComponent {
  v: FlexLayoutView;

  @Input() set view(v: View) {
    this.v = v as FlexLayoutView;
    this.buttonGroup = this.v.config.buttonGroup;
    this.sections = this.v.config.sections;
  }

  get view() {
    return this.v;
  }

  buttonGroup: ButtonGroupView;

  sections: FlexLayoutItem[][];

  identifySection = trackByIndex;

  sectionStyle(item: FlexLayoutItem) {
    return ['height', 'margin'].reduce((previousValue, currentValue) => {
      if (!item[currentValue]) {
        return previousValue;
      }

      return { ...previousValue, [currentValue]: item[currentValue] };
    }, {});
  }
}
