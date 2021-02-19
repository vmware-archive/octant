/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable, OnDestroy } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import {
  Operation,
  PreferencePanel,
  Preferences,
} from '../../models/preference';
import { ThemeService } from '../theme/theme.service';
import { PreferencesEntry } from './preferences.entry';

@Injectable({
  providedIn: 'root',
})
export class PreferencesService implements OnDestroy {
  private electronStore: any;

  public preferences: Map<string, PreferencesEntry<any>> = new Map();
  public preferencesOpened: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
    false
  );

  constructor(private themeService: ThemeService) {
    if (this.isElectron()) {
      const Store = window.require('electron-store');
      this.electronStore = new Store();
    }

    this.preferences.set(
      'navigation.collapsed',
      new PreferencesEntry<boolean>(
        this,
        'navigation.collapsed',
        false,
        'collapsed'
      )
    );

    this.preferences.set(
      'navigation.labels',
      new PreferencesEntry<boolean>(this, 'navigation.labels', true, 'show')
    );

    this.preferences.set(
      'general.pageSize',
      new PreferencesEntry<number>(this, 'general.pageSize', 10, '')
    );

    this.preferences.set(
      'development.frontendUrl',
      new PreferencesEntry<string>(
        this,
        'development.frontendUrl',
        'http://localhost:4200',
        '',
        true
      )
    );

    this.preferences.set(
      'development.embedded',
      new PreferencesEntry<boolean>(
        this,
        'development.embedded',
        true,
        'embedded',
        true
      )
    );

    this.preferences.set(
      'development.verbose',
      new PreferencesEntry<boolean>(this, 'development.verbose', false, 'debug')
    );

    this.preferences.set(
      'theme',
      new PreferencesEntry<string>(
        this,
        'theme',
        this.themeService.themeType.value,
        ''
      )
    );
  }

  ngOnDestroy(): void {
    for (const pref of this.preferences.values()) {
      pref.destroy();
    }
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
    let notificationRequired = false;

    for (const pref of this.preferences.values()) {
      const changed = pref.preferencesChanged(update);
      if (changed && pref.updatesElectron) {
        notificationRequired = true;
      }
    }

    this.updateTheme();

    if (this.isElectron() && notificationRequired) {
      const ipcRenderer = window.require('electron').ipcRenderer;
      ipcRenderer.send('preferences', 'changed');
    }
  }

  public getPreferences(): Preferences {
    const panels: PreferencePanel[] = this.isElectron()
      ? [
          this.getGeneralPanels(),
          this.getNavigationPanels(),
          this.getDeveloperPanels(),
        ]
      : [this.getGeneralPanels(), this.getNavigationPanels()];

    return {
      updateName: 'generalPreferences',
      panels,
    };
  }

  updateTheme() {
    if (
      this.themeService.themeType.value !==
      this.preferences.get('theme').subject.value
    ) {
      this.themeService.switchTheme();
    }
  }

  reset() {
    for (const pref of this.preferences.values()) {
      pref.setDefaultValue();
    }
    this.updateTheme();
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
              value: this.preferences.get('development.embedded').subject.value
                ? 'embedded'
                : 'proxied',
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
              value: this.preferences.get('development.frontendUrl').subject
                .value,
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
              value: this.preferences.get('development.verbose').subject.value
                ? 'debug'
                : 'normal',
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
    const pageSize = this.preferences.get('general.pageSize').subject.value;

    return {
      name: 'General',
      sections: [
        {
          name: 'Color Theme',
          elements: [
            {
              name: 'theme',
              type: 'radio',
              value: this.preferences.get('theme').subject.value,
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
          name: 'Table Pagination',
          elements: [
            {
              name: 'general.pageSize',
              value: pageSize,
              label: 'Page size:',
              type: 'dropdown',
              metadata: {
                type: 'dropdown',
                title: [
                  {
                    metadata: {
                      type: 'text',
                    },
                    config: {
                      value: pageSize,
                    },
                  },
                ],
              },
              config: {
                type: 'label',
                selection: pageSize,
                useSelection: true,
                items: [
                  {
                    name: '10',
                    type: 'text',
                    label: '10',
                  },
                  {
                    name: '20',
                    type: 'text',
                    label: '20',
                  },
                  {
                    name: '50',
                    type: 'text',
                    label: '50',
                  },
                  {
                    name: '100',
                    type: 'text',
                    label: '100',
                  },
                ],
              },
            },
          ],
        },
      ],
    };
  }

  private getNavigationPanels(): PreferencePanel {
    return {
      name: 'Navigation',
      sections: [
        {
          name: 'Navigation',
          elements: [
            {
              name: 'navigation.collapsed',
              type: 'radio',
              value: this.preferences.get('navigation.collapsed').subject.value
                ? 'collapsed'
                : 'expanded',
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
              name: 'navigation.labels',
              type: 'radio',
              value: this.preferences.get('navigation.labels').subject.value
                ? 'show'
                : 'hide',
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
