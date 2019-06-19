// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import trackByIndex from './trackByIndex';
import trackByIdentity from './trackByIdentity';

describe('trackByIndex', () => {
  it('should return the same value as the first param', () => {
    expect(trackByIndex(0)).toBe(0);
    expect(trackByIndex(7)).toBe(7);
  });
});

describe('trackByIdentity', () => {
  it('should return the same value as the second param', () => {
    expect(trackByIdentity(0, 2)).toBe(2);
    expect(trackByIdentity(9, true)).toBe(true);
  });
});
