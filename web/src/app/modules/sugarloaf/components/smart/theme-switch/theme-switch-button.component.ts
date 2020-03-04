// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnInit, Renderer2 } from '@angular/core';
import {
  darkTheme,
  defaultTheme,
  lightTheme,
  Theme,
  ThemeService,
  ThemeType,
} from './theme-switch.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Component({
  selector: 'app-theme-switch-button',
  templateUrl: './theme-switch-button.component.html',
  styleUrls: ['./theme-switch-button.component.scss'],
  providers: [ThemeService, MonacoProviderService],
})
export class ThemeSwitchButtonComponent implements OnInit {
  themeType: ThemeType;

  constructor(
    private themeService: ThemeService,
    private monacoService: MonacoProviderService,
    private renderer: Renderer2
  ) {}

  ngOnInit() {
    this.themeType = this.themeService.currentType();
    this.loadTheme();
  }

  isLightThemeEnabled(): boolean {
    // TODO: this should be in the theme service.
    return this.themeType === 'light';
  }

  loadTheme() {
    // TODO: this should be in the theme service.
    const theme: Theme = this.isLightThemeEnabled() ? lightTheme : darkTheme;

    this.themeService.loadCSS(theme.assetPath);

    [darkTheme, lightTheme].forEach(t =>
      this.renderer.removeClass(document.body, t.type)
    );
    this.renderer.addClass(document.body, theme.type);
    if (this.isLightThemeEnabled()) {
      this.monacoService.changeTheme('vs');
    } else {
      this.monacoService.changeTheme('vs-dark');
    }
  }

  switchTheme() {
    // TODO: this should be in the theme service.
    if (this.isLightThemeEnabled()) {
      this.themeType = 'dark';
      localStorage.setItem('theme', 'dark');
      this.monacoService.changeTheme('vs-dark');
    } else {
      this.themeType = 'light';
      localStorage.setItem('theme', 'light');
      this.monacoService.changeTheme('vs');
    }

    this.loadTheme();
  }
}
