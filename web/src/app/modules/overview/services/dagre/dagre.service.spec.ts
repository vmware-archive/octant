// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import { DagreService } from './dagre.service';

describe('DagreService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: DagreService = TestBed.get(DagreService);
    expect(service).toBeTruthy();
  });
});
