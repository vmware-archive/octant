// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LinkView, View } from '../../../models/content';

@Component({
  selector: 'app-breadcrumb',
  templateUrl: './breadcrumb.component.html',
  styleUrls: ['./breadcrumb.component.scss'],
})
export class BreadcrumbComponent implements OnChanges {
  @Input() path: View[];
  header: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.path.currentValue) {
      const currentPath = changes.path.currentValue as View[];
      const last: LinkView = currentPath[currentPath.length - 1] as LinkView;

      this.header = currentPath.length > 0 ? last.config.value : '';
    }
  }

  identifyPath(index: number, item: View) {
    return `${item.metadata.title}-${index}`;
  }
}
