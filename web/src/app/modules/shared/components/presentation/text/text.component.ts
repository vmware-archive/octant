// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  ChangeDetectorRef,
  Component,
  OnInit,
  SecurityContext,
} from '@angular/core';
import '@cds/core/button/register';
import { ClarityIcons, clipboardIcon } from '@cds/core/icon';
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
  clipboardValue: string;
  copied: boolean;

  isMarkdown: boolean;

  hasStatus = false;

  constructor(
    private readonly sanitizer: DomSanitizer,
    private cdr: ChangeDetectorRef
  ) {
    super();
    ClarityIcons.addIcons(clipboardIcon);
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

    if (view.config.clipboardValue) {
      this.clipboardValue = view.config.clipboardValue;
    }
  }

  copyToClipboard(): void {
    document.addEventListener('copy', (e: ClipboardEvent) => {
      e.clipboardData.setData('text/plain', this.clipboardValue);
      e.preventDefault();
      document.removeEventListener('copy', null);
    });
    document.execCommand('copy');
    this.copied = !this.copied;
    setTimeout(() => {
      this.copied = false;
      this.cdr.detectChanges();
    }, 1500);
  }
}
