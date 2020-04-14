// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ErrorView, View } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-error',
  templateUrl: './error.component.html',
  styleUrls: ['./error.component.scss'],
})
export class ErrorComponent implements OnChanges {
  private v: ErrorView;

  @Input() set view(v: View) {
    this.v = v as ErrorView;
  }
  get view() {
    return this.v;
  }

  source: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as ErrorView;
      this.source = view.config.data;
    }
  }
}
