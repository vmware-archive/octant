/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';
import { LoadingService } from './loading.service';

describe('LoadingService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [LoadingService],
    })
  );

  it('should be created', () => {
    const service: LoadingService = TestBed.get(LoadingService);
    expect(service).toBeTruthy();
  });
});
