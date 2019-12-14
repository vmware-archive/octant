// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { Subject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class SliderService {
  setHeight$: Observable<any>;
  private height = new Subject<number>();

  constructor() {
    this.setHeight$ = this.height.asObservable();
  }

  setHeight(height: number) {
    this.height.next(height);
  }
}
