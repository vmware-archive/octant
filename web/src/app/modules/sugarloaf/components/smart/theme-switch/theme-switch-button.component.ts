// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnInit, Renderer2 } from '@angular/core';
import { ThemeService, ThemeType } from './theme-switch.service';
import { MonacoProviderService } from 'ng-monaco-editor';

@Component({
  selector: 'app-theme-switch-button',
  templateUrl: './theme-switch-button.component.html',
  styleUrls: ['./theme-switch-button.component.scss'],
  providers: [ThemeService, MonacoProviderService],
})
export class ThemeSwitchButtonComponent implements OnInit {
  themeType: ThemeType;

  lightThemeEnabled: boolean;

  constructor(
    private themeService: ThemeService,
    private monacoService: MonacoProviderService,
    private renderer: Renderer2
  ) {}

  ngOnInit() {
    this.themeType = this.themeService.currentType();
    this.lightThemeEnabled = this.themeService.isLightThemeEnabled();
  }

  switchTheme() {
    this.lightThemeEnabled = !this.lightThemeEnabled;
    this.themeService.switchTheme(this.monacoService, this.renderer);
  }
}
