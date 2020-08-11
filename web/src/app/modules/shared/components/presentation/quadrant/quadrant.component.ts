// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import {
  QuadrantValue,
  QuadrantView,
} from 'src/app/modules/shared/models/content';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

const emptyQuadrantValue = { value: '', label: '' };

@Component({
  selector: 'app-view-quadrant',
  templateUrl: './quadrant.component.html',
  styleUrls: ['./quadrant.component.scss'],
})
export class QuadrantComponent extends AbstractViewComponent<QuadrantView> {
  title: string;
  nw: QuadrantValue = emptyQuadrantValue;
  ne: QuadrantValue = emptyQuadrantValue;
  sw: QuadrantValue = emptyQuadrantValue;
  se: QuadrantValue = emptyQuadrantValue;

  constructor(private viewService: ViewService) {
    super();
  }

  update() {
    const view = this.v;
    this.title = this.viewService.viewTitleAsText(view);
    this.nw = view.config.nw;
    this.ne = view.config.ne;
    this.sw = view.config.sw;
    this.se = view.config.se;
  }
}
