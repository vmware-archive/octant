/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { Component, Input, OnInit } from '@angular/core';
import { TitleView } from '../../../models/content';

@Component({
  selector: 'app-view-title',
  templateUrl: './title.component.html',
  styleUrls: ['./title.component.scss'],
})
export class TitleComponent implements OnInit {
  @Input() views: TitleView[];

  constructor() {}

  ngOnInit() {}

  trackBy(index, item) {
    return index;
  }
}
