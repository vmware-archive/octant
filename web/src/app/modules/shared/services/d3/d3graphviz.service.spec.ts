// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import { D3GraphvizService } from './d3graphviz.service';

describe('DagreService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: D3GraphvizService = TestBed.get(D3GraphvizService);
    expect(service).toBeTruthy();
  });
});
