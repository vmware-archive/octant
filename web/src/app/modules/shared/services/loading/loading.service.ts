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
  takeUntil,
} from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class LoadingService {
  public requestComplete = new BehaviorSubject<boolean>(false);

  constructor() {}

  withDelay(
    watch: Observable<boolean>,
    after: number,
    atLeast: number
  ): Observable<boolean> {
    const loadingTimer = timer(after).pipe(takeUntil(watch));
    const holdTimer = timer(after + atLeast);

    return merge<boolean>(
      loadingTimer.pipe(mapTo(true)),
      combineLatest([watch, holdTimer]).pipe(mapTo(false))
    ).pipe(startWith(false), distinctUntilChanged());
  }
}
