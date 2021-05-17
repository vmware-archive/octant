// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { ThemeService } from '../../../../shared/services/theme/theme.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-theme-switch-button',
  templateUrl: './theme-switch-button.component.html',
  styleUrls: ['./theme-switch-button.component.scss'],
})
export class ThemeSwitchButtonComponent implements OnInit, OnDestroy {
  @Input() public collapsed: boolean;

  private subscriptionTheme: Subscription;
  lightThemeEnabled: boolean;

  private onThemeChange: () => void;

  constructor(private themeService: ThemeService) {
    // we want a new instance of the handler for each component instance
    this.onThemeChange = () => {
      this.lightThemeEnabled = this.themeService.isLightThemeEnabled();
    };
    this.onThemeChange();
  }

  ngOnInit() {
    this.subscriptionTheme = this.themeService.themeType.subscribe(() =>
      this.onThemeChange()
    );
  }

  ngOnDestroy() {
    if (this.subscriptionTheme) {
      this.subscriptionTheme.unsubscribe();
    }
  }

  switchTheme(): void {
    this.themeService.switchTheme();
  }
}
