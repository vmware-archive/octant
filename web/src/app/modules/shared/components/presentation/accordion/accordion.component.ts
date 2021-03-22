// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnInit } from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import { AccordionRow, AccordionView } from '../../../models/content';
import { ViewService } from '../../../services/view/view.service';

@Component({
  selector: 'app-view-accordion',
  templateUrl: './accordion.component.html',
  styleUrls: ['./accordion.component.scss'],
})
export class AccordionComponent
  extends AbstractViewComponent<AccordionView>
  implements OnInit {
  rows: AccordionRow[];
  title: string;
  allowMultipleExpanded: boolean;
  trackByIndex = trackByIndex;

  constructor(private viewService: ViewService) {
    super();
  }

  update() {
    const view = this.v;
    this.title = this.viewService.viewTitleAsText(view);
    this.rows = view.config.rows;
    this.allowMultipleExpanded = view.config.allowMultipleExpanded;
  }
}
