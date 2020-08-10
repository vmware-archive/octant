// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnInit } from '@angular/core';
import { ThemeService } from './theme-switch.service';

@Component({
  selector: 'app-theme-switch-button',
  templateUrl: './theme-switch-button.component.html',
  styleUrls: ['./theme-switch-button.component.scss'],
  providers: [ThemeService],
})
export class ThemeSwitchButtonComponent implements OnInit {
  @Input() public collapsed: boolean;

  lightThemeEnabled: boolean;

  constructor(private themeService: ThemeService) {}

  ngOnInit() {
    this.lightThemeEnabled = this.themeService.isLightThemeEnabled();
  }

  switchTheme(): Promise<any> {
    return this.themeService
      .switchTheme()
      .then(() => {
        this.lightThemeEnabled = this.themeService.isLightThemeEnabled();
      })
      .catch(e => {
        console.error('Unable to switch theme:', e);
      });
  }
}
