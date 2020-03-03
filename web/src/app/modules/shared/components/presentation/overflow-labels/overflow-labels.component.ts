// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit } from '@angular/core';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { LabelFilterService } from '../../../services/label-filter/label-filter.service';

interface Labels {
  [key: string]: string;
}

@Component({
  selector: 'app-overflow-labels',
  templateUrl: './overflow-labels.component.html',
  styleUrls: ['./overflow-labels.component.scss'],
})
export class OverflowLabelsComponent implements OnInit {
  @Input() labels: Labels;
  @Input() numberShownLabels = 2;

  showLabels: Labels[];
  overflowLabels: Labels[];
  trackByIdentity = trackByIdentity;

  ngOnInit() {
    const labelsEntries = Object.entries({ ...this.labels });

    if (this.numberShownLabels <= labelsEntries.length) {
      this.showLabels = labelsEntries
        .slice(0, this.numberShownLabels)
        .map(label => ({ [label[0]]: label[1] }));

      this.overflowLabels = labelsEntries
        .slice(this.numberShownLabels)
        .map(label => ({ [label[0]]: label[1] }));
    } else {
      this.showLabels = labelsEntries.map(label => ({ [label[0]]: label[1] }));
    }
  }

  filterLabel(key: string, value: string) {
    this.labelFilter.add({ key, value });
  }

  constructor(private labelFilter: LabelFilterService) {}
}
