// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LoadingView, View } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-loading',
  templateUrl: './loading.component.html',
  styleUrls: ['./loading.component.scss'],
})
export class LoadingComponent implements OnChanges {
  private v: LoadingView;

  @Input() set view(v: View) {
    this.v = v as LoadingView;
  }
  get view() {
    return this.v;
  }

  value: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LoadingView;
      this.value = view.config.value;
    }
  }
}
