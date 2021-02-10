/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject, Subscription } from 'rxjs';
import {
  Operation,
  PreferencePanel,
  Preferences,
} from '../../models/preference';
import { ThemeService } from '../theme/theme.service';
import { skip } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class PreferencesService implements OnDestroy {
  private subscriptionCollapsed: Subscription;
  private subscriptionLabels: Subscription;
  private subscriptionTheme: Subscription;
  private subscriptionFrontendUrl: Subscription;
  private subscriptionVerbose: Subscription;
  private subscriptionEmbedded: Subscription;
  private electronStore: any;

  public preferencesOpened: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
    false
  );

  public navCollapsed: BehaviorSubject<boolean>;
  public showLabels: BehaviorSubject<boolean>;
  public frontendUrl: BehaviorSubject<string>;
  public verbose: BehaviorSubject<boolean>;
  public embedded: BehaviorSubject<boolean>;

  constructor(private themeService: ThemeService) {
    if (this.isElectron()) {
      const Store = window.require('electron-store');
      this.electronStore = new Store();
    }

    this.navCollapsed = new BehaviorSubject<boolean>(
      JSON.parse(this.getStoredValue('navigation.collapsed', false))
    );
    this.showLabels = new BehaviorSubject<boolean>(
      JSON.parse(this.getStoredValue('navigation.labels', true))
    );
    this.frontendUrl = new BehaviorSubject<string>(
      this.getStoredValue('development.frontendUrl', 'http://localhost:4200')
    );

    this.verbose = new BehaviorSubject<boolean>(
      this.getStoredValue('development.verbose', false)
    );

    this.embedded = new BehaviorSubject<boolean>(
      this.getStoredValue('development.embedded', true)
    );

    this.subscriptionCollapsed = this.navCollapsed.subscribe(col => {
      this.setStoredValue('navigation.collapsed', col);
    });

    this.subscriptionLabels = this.showLabels.subscribe(labels => {
      this.setStoredValue('navigation.labels', labels);
    });

    this.subscriptionFrontendUrl = this.frontendUrl.subscribe(url => {
      this.setStoredValue('development.frontendUrl', url);
    });

    this.subscriptionVerbose = this.verbose.subscribe(verbose => {
      this.setStoredValue('development.verbose', verbose);
    });

    this.subscriptionEmbedded = this.embedded.subscribe(embedded => {
      this.setStoredValue('development.embedded', embedded);
    });

    this.subscriptionTheme = this.themeService.themeType
      .pipe(skip(1))
      .subscribe(theme => {
        this.setStoredValue('theme', theme);
      });

    this.themeService.themeType.next(
      this.getStoredValue('theme', this.themeService.themeType.value)
    );
  }

  ngOnDestroy(): void {
    this.subscriptionCollapsed?.unsubscribe();
    this.subscriptionLabels?.unsubscribe();
    this.subscriptionTheme?.unsubscribe();
    this.subscriptionFrontendUrl?.unsubscribe();
    this.subscriptionVerbose?.unsubscribe();
    this.subscriptionEmbedded?.unsubscribe();
  }

  setStoredValue(key: string, value: any) {
    if (this.isElectron()) {
      this.electronStore.set(key, value);
    } else {
      localStorage.setItem(key, value);
    }
  }

  getStoredValue(key: string, defaultValue: any) {
    if (this.isElectron()) {
      return this.electronStore.get(key, defaultValue);
    } else {
      return localStorage.getItem(key) || defaultValue;
    }
  }

  isElectron(): boolean {
    if (typeof process === 'undefined') {
      return false;
    }
    return (
      process && process.versions && process.versions.electron !== undefined
    );
  }

  public preferencesChanged(update: Preferences) {
    const collapsed = update['general.navigation'] === 'collapsed';
    const showLabels = update['general.labels'] === 'show';
    const isLightTheme = update['general.theme'] === 'light';
    const frontendUrl = update['development.frontendUrl'];
    const verbose = update['development.verbose'] === 'debug';
    const embedded = update['development.embedded'] === 'embedded';
    let notificationRequired = false;

    if (this.showLabels.value !== showLabels) {
      this.showLabels.next(showLabels);
    }

    if (this.navCollapsed.value !== collapsed) {
      this.navCollapsed.next(collapsed);
    }

    if (this.themeService.isLightThemeEnabled() !== isLightTheme) {
      this.themeService.switchTheme();
    }

    if (this.frontendUrl.value !== frontendUrl) {
      notificationRequired = true;
      this.frontendUrl.next(frontendUrl);
    }

    if (this.embedded.value !== embedded) {
      notificationRequired = true;
      this.embedded.next(embedded);
    }

    if (this.verbose.value !== verbose) {
      this.verbose.next(verbose);
    }

    if (this.isElectron() && notificationRequired) {
      const ipcRenderer = window.require('electron').ipcRenderer;
      ipcRenderer.send('preferences', 'changed');
    }
  }

  public getPreferences(): Preferences {
    const panels: PreferencePanel[] = this.isElectron()
      ? [this.getGeneralPanels(), this.getDeveloperPanels()]
      : [this.getGeneralPanels()];

    return {
      updateName: 'generalPreferences',
      panels,
    };
  }

  private getDeveloperPanels(): PreferencePanel {
    return {
      name: 'Development',
      sections: [
        {
          name: 'Frontend Source',
          elements: [
            {
              name: 'development.embedded',
              type: 'radio',
              value: this.embedded.value ? 'embedded' : 'proxied',
              config: {
                values: [
                  {
                    label: 'Embedded',
                    value: 'embedded',
                  },
                  {
                    label: 'Proxied',
                    value: 'proxied',
                  },
                ],
              },
            },
            {
              name: 'development.frontendUrl',
              type: 'text',
              value: this.frontendUrl.value,
              disableConditions: [
                {
                  lhs: 'development.embedded',
                  op: Operation.Equal,
                  rhs: 'proxied',
                },
              ],
              config: {
                label: 'Frontend Proxy Url',
                placeholder: 'http://example.com',
              },
            },
          ],
        },
        {
          name: 'Logging verbosity (requires restart)',
          elements: [
            {
              name: 'development.verbose',
              type: 'radio',
              value: this.verbose.value ? 'debug' : 'normal',
              config: {
                values: [
                  {
                    label: 'Debug',
                    value: 'debug',
                  },
                  {
                    label: 'Normal',
                    value: 'normal',
                  },
                ],
              },
            },
          ],
        },
      ],
    };
  }

  private getGeneralPanels(): PreferencePanel {
    return {
      name: 'General',
      sections: [
        {
          name: 'Color Theme',
          elements: [
            {
              name: 'general.theme',
              type: 'radio',
              value: this.themeService.isLightThemeEnabled() ? 'light' : 'dark',
              config: {
                values: [
                  {
                    label: 'Dark',
                    value: 'dark',
                  },
                  {
                    label: 'Light',
                    value: 'light',
                  },
                ],
              },
            },
          ],
        },
        {
          name: 'Navigation',
          elements: [
            {
              name: 'general.navigation',
              type: 'radio',
              value: this.navCollapsed.value ? 'collapsed' : 'expanded',
              config: {
                values: [
                  {
                    label: 'Expanded',
                    value: 'expanded',
                  },
                  {
                    label: 'Collapsed',
                    value: 'collapsed',
                  },
                ],
              },
            },
          ],
        },
        {
          name: 'Navigation labels',
          elements: [
            {
              name: 'general.labels',
              type: 'radio',
              value: this.showLabels.value ? 'show' : 'hide',
              config: {
                values: [
                  {
                    label: 'Show Labels',
                    value: 'show',
                  },
                  {
                    label: 'Hide Labels',
                    value: 'hide',
                  },
                ],
              },
            },
          ],
        },
      ],
    };
  }
}
