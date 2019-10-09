// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import { PodLogsService } from './pod-logs.service';

describe('PodLogsService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: PodLogsService = TestBed.get(PodLogsService);
    expect(service).toBeTruthy();
  });
});
