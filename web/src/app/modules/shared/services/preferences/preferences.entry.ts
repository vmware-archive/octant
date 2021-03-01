/*
 * Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { BehaviorSubject, Subscription } from 'rxjs';
import { PreferencesService } from './preferences.service';
import { Preferences } from '../../models/preference';

export class PreferencesEntry<T> {
  private subscription: Subscription;
  public subject: BehaviorSubject<T>;

  constructor(
    private preferencesService: PreferencesService,
    public id: string,
    private defaultValue: T,
    private defaultText: string,
    public updatesElectron: boolean = false
  ) {
    if (typeof this.defaultValue !== 'string') {
      this.subject = new BehaviorSubject<T>(
        JSON.parse(
          preferencesService.getStoredValue(this.id, this.defaultValue)
        )
      );
    } else {
      this.subject = new BehaviorSubject<T>(
        preferencesService.getStoredValue(this.id, this.defaultValue)
      );
    }

    this.subscription = this.subject.subscribe(val => {
      preferencesService.setStoredValue(this.id, val);
    });
  }

  public preferencesChanged(update: Preferences) {
    switch (typeof this.defaultValue) {
      case 'boolean':
        const val = (update[this.id] === this.defaultText) as unknown;
        if (this.subject.value !== (val as T)) {
          this.subject.next(val as T);
          return true;
        }
        break;
      default:
        const newValue = update[this.id];
        if (newValue && this.subject.value !== newValue) {
          this.subject.next(newValue);
          return true;
        }
        break;
    }
    return false;
  }

  public setDefaultValue() {
    this.subject.next(this.defaultValue);
  }

  public destroy() {
    this.subscription?.unsubscribe();
  }
}
