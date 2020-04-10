// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import { PathItem } from '../../../models/content';

@Component({
  selector: 'app-breadcrumb',
  templateUrl: './breadcrumb.component.html',
  styleUrls: ['./breadcrumb.component.scss'],
})
export class BreadcrumbComponent {
  @Input() path: PathItem[];

  constructor() {}

  identifyPath(index: number, item: PathItem) {
    return `${item.title}-${index}`;
  }
}
