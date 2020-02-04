// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Inject, Injectable } from '@angular/core';
import { DOCUMENT } from '@angular/common';

export type ThemeType = 'light' | 'dark';

export interface Theme {
  type: ThemeType;
  assetPath: string;
}

/**
 * Dark theme
 */
export const darkTheme: Theme = {
  type: 'dark',
  assetPath: 'assets/css/clr-ui-dark.min.css',
};

/**
 * Light theme
 */
export const lightTheme: Theme = {
  type: 'light',
  assetPath: 'assets/css/clr-ui.min.css',
};

export const defaultTheme = lightTheme;

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

  currentType() {
    const themeType = localStorage.getItem('theme') as ThemeType;
    return themeType || defaultTheme.type;
  }
}
