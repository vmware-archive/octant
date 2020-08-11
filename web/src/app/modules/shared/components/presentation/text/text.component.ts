// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component } from '@angular/core';
import { TextView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss'],
})
export class TextComponent extends AbstractViewComponent<TextView> {
  value: string;

  isMarkdown: boolean;

  hasStatus = false;

  constructor() {
    super();
  }

  update() {
    const view = this.v;
    this.value = view.config.value;
    this.isMarkdown = view.config.isMarkdown;

    if (view.config.status) {
      this.hasStatus = true;
    }
  }
}
