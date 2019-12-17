// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable, EventEmitter, Output } from '@angular/core';
import { Subject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class SliderService {
  setHeight$: Observable<any>;
  resetDefault$: Observable<any>;
  resetExpandedDefault$: Observable<any>;
  @Output() resizedSliderEvent = new EventEmitter<any>();

  private height = new Subject<number>();

  constructor() {
    this.setHeight$ = this.height.asObservable();
  }

  setHeight(height: number) {
    this.height.next(height);
  }

  resetDefault() {
    // Approximate conversion from 1.5rem
    this.height.next(36);
  }

  resetExpandedDefault() {
    // Approximate conversion from 12rem
    this.height.next(288);
  }
}
