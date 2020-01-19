/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface NavigatorStatus {
  history: string[];
  index: number;
}

@Injectable({
  providedIn: 'root',
})
export class NavigatorService {
  private internalHistory: string[] = [];
  private internalIndex = -1;

  private source = new BehaviorSubject<NavigatorStatus>(null);
  status = this.source.asObservable();

  constructor() {}

  addHistory(path: string) {
    console.log(
      `tail: ${
        this.internalHistory[this.internalHistory.length - 1]
      }; path: ${path}`
    );

    if (this.internalHistory.length > 1) {
      if (this.internalHistory[this.internalHistory.length - 2] === path) {
        return;
      }
    } else if (this.internalHistory[this.internalHistory.length - 1] === path) {
      return;
    }

    this.internalHistory = [...this.internalHistory, path];
    this.internalIndex++;

    this.source.next({
      history: this.internalHistory,
      index: this.internalIndex,
    });
    console.log({ history: this.internalHistory, index: this.internalIndex });
  }
}
