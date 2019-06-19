// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { AnnotationsView } from 'src/app/models/content';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-view-annotations',
  templateUrl: './annotations.component.html',
  styleUrls: ['./annotations.component.scss'],
})
export class AnnotationsComponent implements OnChanges {
  @Input() view: AnnotationsView;
  annotations: { [key: string]: string };
  annotationKeys: string[];
  trackByIdentity = trackByIdentity;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as AnnotationsView;
      this.annotations = view.config.annotations;
      this.annotationKeys = Object.keys(this.annotations);
    }
  }
}
