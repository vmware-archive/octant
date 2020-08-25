// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import {
  ContainerDef,
  ContainersView,
} from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-containers',
  templateUrl: './containers.component.html',
  styleUrls: ['./containers.component.scss'],
})
export class ContainersComponent extends AbstractViewComponent<ContainersView> {
  containers: ContainerDef[];

  constructor() {
    super();
  }

  update() {
    this.containers = this.v.config.containers;
  }

  trackItem(index: number, item: ContainerDef): string {
    return item.name;
  }
}
