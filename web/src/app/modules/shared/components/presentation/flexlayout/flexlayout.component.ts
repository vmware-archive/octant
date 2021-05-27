// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import {
  ButtonGroupView,
  FlexLayoutItem,
  Alert,
} from 'src/app/modules/shared/models/content';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-flexlayout',
  templateUrl: './flexlayout.component.html',
  styleUrls: ['./flexlayout.component.scss'],
})
export class FlexlayoutComponent extends AbstractViewComponent<any> {
  buttonGroup: ButtonGroupView;
  sections: FlexLayoutItem[][];
  alert: Alert;

  identifySection = trackByIndex;

  constructor() {
    super();
  }

  update() {
    this.buttonGroup = this.v.config.buttonGroup;
    this.sections = this.v.config.sections;
    this.alert = this.v.config.alert;
  }

  sectionStyle(item: FlexLayoutItem) {
    return ['height', 'margin'].reduce((previousValue, currentValue) => {
      if (!item[currentValue]) {
        return previousValue;
      }

      return { ...previousValue, [currentValue]: item[currentValue] };
    }, {});
  }
}
