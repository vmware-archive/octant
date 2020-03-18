// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Inject, Injectable, Renderer2 } from '@angular/core';
import { DOCUMENT } from '@angular/common';
import { MonacoProviderService } from 'ng-monaco-editor';

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
  private themeType: ThemeType;
  private currentTheme: Theme;

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

  loadTheme(monacoService: MonacoProviderService, renderer: Renderer2): void {
    this.currentTheme = this.isLightThemeEnabled() ? lightTheme : darkTheme;
    this.loadCSS(this.currentTheme.assetPath);

    [darkTheme, lightTheme].forEach(t =>
      renderer.removeClass(document.body, t.type)
    );
    renderer.addClass(document.body, this.currentTheme.type);
    if (this.isLightThemeEnabled()) {
      monacoService.changeTheme('vs');
    } else {
      monacoService.changeTheme('vs-dark');
    }
  }

  switchTheme(monacoService: MonacoProviderService, renderer: Renderer2): void {
    if (this.isLightThemeEnabled()) {
      this.themeType = 'dark';
      localStorage.setItem('theme', 'dark');
      monacoService.changeTheme('vs-dark');
    } else {
      this.themeType = 'light';
      localStorage.setItem('theme', 'light');
      monacoService.changeTheme('vs');
    }

    this.loadTheme(monacoService, renderer);
  }

  currentType(): ThemeType {
    this.themeType = localStorage.getItem('theme') as ThemeType;
    return this.themeType || defaultTheme.type;
  }

  isLightThemeEnabled(): boolean {
    this.themeType = localStorage.getItem('theme') as ThemeType;
    return this.themeType === lightTheme.type;
  }
}
