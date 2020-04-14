// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TextView, View } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss'],
})
export class TextComponent implements OnChanges {
  private v: TextView;

  @Input() set view(v: TextView) {
    this.v = v as TextView;
  }

  get view(): TextView {
    return this.v;
  }

  value: string;

  isMarkdown: boolean;

  hasStatus = false;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as TextView;
      this.value = view.config.value;
      this.isMarkdown = view.config.isMarkdown;

      if (view.config.status) {
        this.hasStatus = true;
      }
    }
  }
}
