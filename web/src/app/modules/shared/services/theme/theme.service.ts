// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Inject, Injectable, Renderer2, RendererFactory2 } from '@angular/core';
import { DOCUMENT } from '@angular/common';
import { BehaviorSubject } from 'rxjs';

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

export const defaultTheme = window.matchMedia('(prefers-color-scheme:dark)')
  .matches
  ? darkTheme
  : lightTheme;

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  public themeType: BehaviorSubject<ThemeType> = new BehaviorSubject<ThemeType>(
    defaultTheme.type
  );

  private renderer: Renderer2;

  constructor(
    @Inject(DOCUMENT) private document: Document,
    rendererFactory: RendererFactory2
  ) {
    this.renderer = rendererFactory.createRenderer(null, null);
  }

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

  loadTheme(): void {
    const currentTheme = this.isLightThemeEnabled() ? lightTheme : darkTheme;
    this.loadCSS(currentTheme.assetPath);

    [darkTheme, lightTheme].forEach(t =>
      this.renderer.removeClass(this.document.body, t.type)
    );
    this.renderer.addClass(this.document.body, currentTheme.type);
    this.renderer.setAttribute(
      this.document.body,
      'cds-theme',
      currentTheme.type
    );
  }

  switchTheme(): void {
    const theme = this.isLightThemeEnabled() ? 'dark' : 'light';
    this.themeType.next(theme);
    this.loadTheme();
  }

  isLightThemeEnabled(): boolean {
    return this.themeType.value === lightTheme.type;
  }
}
