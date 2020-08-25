// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { YAMLView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-yaml',
  templateUrl: './yaml.component.html',
  styleUrls: ['./yaml.component.scss'],
})
export class YamlComponent extends AbstractViewComponent<YAMLView> {
  source: string;

  constructor() {
    super();
  }

  update() {
    this.source = this.v.config.data;
  }
}
