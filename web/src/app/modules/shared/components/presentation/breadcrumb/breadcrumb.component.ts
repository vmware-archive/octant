// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { PathItem } from '../../../models/content';

@Component({
  selector: 'app-breadcrumb',
  templateUrl: './breadcrumb.component.html',
  styleUrls: ['./breadcrumb.component.scss'],
})
export class BreadcrumbComponent implements OnChanges {
  @Input() path: PathItem[];
  header: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.path.currentValue) {
      const currentPath = changes.path.currentValue as PathItem[];
      this.header =
        currentPath.length > 0 ? currentPath[currentPath.length - 1].title : '';
    }
  }

  identifyPath(index: number, item: PathItem) {
    return `${item.title}-${index}`;
  }
}
