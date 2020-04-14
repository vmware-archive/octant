// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ContainerDef, ContainersView, View } from '../../../models/content';

@Component({
  selector: 'app-view-containers',
  templateUrl: './containers.component.html',
  styleUrls: ['./containers.component.scss'],
})
export class ContainersComponent implements OnChanges {
  private v: ContainersView;

  @Input() set view(v: View) {
    this.v = v as ContainersView;
  }
  get view() {
    return this.v;
  }

  containers: ContainerDef[];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as ContainersView;
      this.containers = view.config.containers;
    }
  }

  trackItem(index: number, item: ContainerDef): string {
    return item.name;
  }
}
