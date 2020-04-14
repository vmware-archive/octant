/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { CodeView, View } from '../../../models/content';

@Component({
  selector: 'app-view-code',
  templateUrl: './code.component.html',
  styleUrls: ['./code.component.scss'],
})
export class CodeComponent implements OnChanges {
  private v: CodeView;

  @Input() set view(v: View) {
    this.v = v as CodeView;
  }

  get view() {
    return this.v;
  }

  value: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as CodeView;
      this.value = view.config.value;
    }
  }

  copyToClipboard(text) {
    document.addEventListener('copy', (e: ClipboardEvent) => {
      e.clipboardData.setData('text/plain', text);
      e.preventDefault();
      document.removeEventListener('copy', null);
    });
    document.execCommand('copy');
  }
}
