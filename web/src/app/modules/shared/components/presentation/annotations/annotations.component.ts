// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { AnnotationsView, View } from '../../../../shared/models/content';
import trackByIdentity from '../../../../../util/trackBy/trackByIdentity';

@Component({
  selector: 'app-view-annotations',
  templateUrl: './annotations.component.html',
  styleUrls: ['./annotations.component.scss'],
})
export class AnnotationsComponent implements OnChanges {
  private v: AnnotationsView;

  @Input() set view(v: View) {
    this.v = v as AnnotationsView;
  }

  get view() {
    return this.v;
  }

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
