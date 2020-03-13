/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { EditorView, View } from 'src/app/modules/shared/models/content';
import { ThemeService } from '../../../../sugarloaf/components/smart/theme-switch/theme-switch.service';

interface Options {
  readOnly: boolean;
  language: string;
  theme: string;
}

@Component({
  selector: 'app-view-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss'],
})
export class EditorComponent implements OnChanges {
  private v: EditorView;

  @Input() set view(v: View) {
    this.v = v as EditorView;
  }

  get view() {
    return this.v;
  }

  value: string;
  options: Options;

  constructor(private themeService: ThemeService) {
    this.options = {} as Options;
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as EditorView;
      this.value = view.config.value;
      this.options.readOnly = view.config.readOnly;
      this.options.language = view.config.language;
      this.options.theme =
        this.themeService.currentType() === 'dark' ? 'vs-dark' : 'vs';
    }
  }
}
