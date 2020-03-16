/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { NgZone } from '@angular/core';
import { RelativePipe } from './relative.pipe';

class NgZoneMock {
  runOutsideAngular(fn: () => void) {
    return fn();
  }
  run(fn: () => void) {
    return fn();
  }
}

describe('RelativePipe', () => {
  const pipe = new RelativePipe(null, new NgZoneMock() as NgZone);
  const now = new Date(1583971407000);

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('Transform to 0s', () => {
    expect(pipe.transform(1583971407, now)).toEqual('0s');
  });

  it('Transform to 1m', () => {
    expect(pipe.transform(1583971346, now)).toEqual('1m');
  });

  it('Transform to 1h', () => {
    expect(pipe.transform(1583967806, now)).toEqual('1h');
  });

  it('Transform to 1d', () => {
    expect(pipe.transform(1583885006, now)).toEqual('1d');
  });
});
