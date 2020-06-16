/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { BehaviorSubject, combineLatest, merge, Observable, timer } from 'rxjs';
import {
  distinctUntilChanged,
  mapTo,
  startWith,
  takeWhile,
} from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class LoadingService {
  public requestComplete = new BehaviorSubject<boolean>(false);

  constructor() {}

  public showSpinner: Observable<boolean> = merge(
    timer(650).pipe(
      // show only if operation > 650ms
      mapTo(true),
      takeWhile(() => !this.requestComplete.value)
    ), // if shown, stay at least 1sec to prevent flicker
    combineLatest([this.requestComplete, timer(1650)]).pipe(mapTo(false))
  ).pipe(startWith(false), distinctUntilChanged());
}
