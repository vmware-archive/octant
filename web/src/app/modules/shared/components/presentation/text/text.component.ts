// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnInit, SecurityContext } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { TextView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { parse } from 'marked';

@Component({
  selector: 'app-view-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss'],
})
export class TextComponent
  extends AbstractViewComponent<TextView>
  implements OnInit {
  value: string | SafeHtml;

  isMarkdown: boolean;

  hasStatus = false;

  constructor(private readonly sanitizer: DomSanitizer) {
    super();
  }

  update() {
    const view = this.v;

    this.isMarkdown = view.config.isMarkdown;

    if (view.config.isMarkdown) {
      const html = parse(view.config.value);
      this.value = view.config.trustedContent
        ? this.sanitizer.bypassSecurityTrustHtml(html)
        : this.sanitizer.sanitize(SecurityContext.HTML, html);
    } else {
      this.value = view.config.value;
    }

    if (view.config.status) {
      this.hasStatus = true;
    }
  }
}
