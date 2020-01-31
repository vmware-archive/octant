// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LabelsView, View } from 'src/app/shared/models/content';
import { LabelFilterService } from 'src/app/shared/services/label-filter/label-filter.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-view-labels',
  templateUrl: './labels.component.html',
  styleUrls: ['./labels.component.scss'],
})
export class LabelsComponent implements OnChanges {
  private v: LabelsView;

  @Input() set view(v: View) {
    this.v = v as LabelsView;
  }
  get view() {
    return this.v;
  }

  title: View[];
  labelKeys: string[];
  labels: { [key: string]: string };
  trackByIdentity = trackByIdentity;

  constructor(private labelFilter: LabelFilterService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LabelsView;

      this.title = view.metadata.title;
      this.labels = view.config.labels;
      this.labelKeys = Object.keys(this.labels);
    }
  }

  click(key: string, value: string) {
    this.labelFilter.add({ key, value });
  }
}
