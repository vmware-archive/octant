// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { AnnotationsView } from 'src/app/modules/shared/models/content';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-annotations',
  templateUrl: './annotations.component.html',
  styleUrls: ['./annotations.component.scss'],
})
export class AnnotationsComponent extends AbstractViewComponent<
  AnnotationsView
> {
  annotations: { [key: string]: string };
  annotationKeys: string[];
  trackByIdentity = trackByIdentity;

  constructor() {
    super();
  }

  update() {
    const view = this.v;
    this.annotations = view.config.annotations;
    this.annotationKeys = this.annotations ? Object.keys(this.annotations) : [];
  }
}
