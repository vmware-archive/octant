// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Injectable, Inject } from '@angular/core';
import { DOCUMENT } from '@angular/common';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  constructor(@Inject(DOCUMENT) private document: Document) {}

  loadCSS(route: string) {
    const head = this.document.getElementsByTagName('head')[0];
    const themeLink = this.document.getElementById(
      'client-theme'
    ) as HTMLLinkElement;

    if (themeLink) {
      themeLink.href = route;
    } else {
      const style = this.document.createElement('link');
      style.id = 'client-theme';
      style.rel = 'stylesheet';
      style.href = `${route}`;

      head.appendChild(style);
    }
  }
}
