// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Inject,
  Injectable,
  Renderer2,
  RendererFactory2,
  OnDestroy,
} from '@angular/core';
import { DOCUMENT } from '@angular/common';
import { EventEmitter } from 'events';

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
export class ThemeService implements OnDestroy {
  private emitter: EventEmitter;
  private themeType: ThemeType;
  private renderer: Renderer2;
  private storageEventHandler: (e: StorageEvent) => void;

  constructor(
    @Inject(DOCUMENT) private document: Document,
    rendererFactory: RendererFactory2
  ) {
    this.emitter = new EventEmitter();

    const themeType = localStorage.getItem('theme') as ThemeType;
    this.themeType = themeType || defaultTheme.type;
    
    this.renderer = rendererFactory.createRenderer(null, null);

    this.storageEventHandler = (e: StorageEvent): void => {
      if (e.key === 'theme' && e.newValue !== this.themeType) {
        // another window switched the theme
        this.switchTheme();
      }
    };
    addEventListener('storage', this.storageEventHandler);
  }

  ngOnDestroy(): void {
    removeEventListener('storage', this.storageEventHandler);
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
  }

  switchTheme(): void {
    this.themeType = this.isLightThemeEnabled() ? 'dark' : 'light';
    localStorage.setItem('theme', this.themeType);
    this.loadTheme();

    this.emitter.emit('change');
  }

  isLightThemeEnabled(): boolean {
    return this.themeType === lightTheme.type;
  }

  onChange(listener: () => void): void {
    this.emitter.on('change', listener);
  }

  offChange(listener: () => void): void {
    this.emitter.off('change', listener);
  }
}
