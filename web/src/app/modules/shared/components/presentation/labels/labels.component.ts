// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import { LabelsView, View } from 'src/app/modules/shared/models/content';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-labels',
  templateUrl: './labels.component.html',
  styleUrls: ['./labels.component.scss'],
})
export class LabelsComponent extends AbstractViewComponent<LabelsView> {
  title: string;
  labelKeys: string[];
  labels: { [key: string]: string };
  trackByIdentity = trackByIdentity;

  constructor(private viewService: ViewService) {
    super();
  }

  update() {
    const view = this.v;
    this.title = this.viewService.viewTitleAsText(view);
    this.labels = view.config.labels;
  }
}
