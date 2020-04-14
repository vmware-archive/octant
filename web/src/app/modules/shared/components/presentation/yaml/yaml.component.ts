// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View, YAMLView } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-yaml',
  templateUrl: './yaml.component.html',
  styleUrls: ['./yaml.component.scss'],
})
export class YamlComponent implements OnChanges {
  private v: YAMLView;

  @Input() set view(v: View) {
    this.v = v as YAMLView;
  }
  get view() {
    return this.v;
  }

  source: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as YAMLView;
      this.source = view.config.data;
    }
  }
}
