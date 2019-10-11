// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnInit, Renderer2 } from '@angular/core';
import { ThemeService } from './theme-switch.service';

type Theme = 'light' | 'dark'

@Component({
  selector: 'app-theme-switch-button',
  templateUrl: './theme-switch-button.component.html',
  styleUrls: ['./theme-switch-button.component.scss'],
  providers: [ThemeService]
})
export class ThemeSwitchButtonComponent implements OnInit {
  theme: Theme;

  constructor(
    private themeService: ThemeService,
    private renderer: Renderer2,
  ) { }

  ngOnInit() {
    this.theme = localStorage.getItem('theme') as Theme || 'light';
    this.loadTheme();
  }

  isLightThemeEnabled(): boolean {
    return this.theme === 'light'
  }

  loadTheme() {
    this.themeService.loadCSS(
      this.isLightThemeEnabled() ? 'assets/css/clr-ui.min.css' : 'assets/css/clr-ui-dark.min.css'
    );
    this.renderer.removeClass(document.body, this.isLightThemeEnabled() ? 'dark' : 'light');
    this.renderer.addClass(document.body, this.isLightThemeEnabled() ? 'light' : 'dark');
  }

  switchTheme() {
    if (this.isLightThemeEnabled()) {
      this.theme = 'dark';
      localStorage.setItem('theme', 'dark');
    } else {
      this.theme = 'light';
      localStorage.setItem('theme', 'light');
    }

    this.loadTheme();
  }
}
