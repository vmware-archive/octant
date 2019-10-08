// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TextView } from 'src/app/models/content';

@Component({
  selector: 'app-view-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss'],
})
export class TextComponent implements OnChanges {
  @Input() view: TextView;

  value: string;

  isMarkdown: boolean;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as TextView;
      this.value = view.config.value;
      this.isMarkdown = view.config.isMarkdown;
    }
  }
}
