// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';
import { SliderService } from './slider.service';

describe('SliderService', () => {
  let service: SliderService;

  beforeEach(() => {
    service = TestBed.get(SliderService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('set value', () => {
    service.setHeight(100);
    service.setHeight$.subscribe(current => expect(current).toEqual(100));
  });

  it('reset to default', () => {
    service.resetDefault();
    service.setHeight$.subscribe(current => expect(current).toEqual(36));
  });
});
