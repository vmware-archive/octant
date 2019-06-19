// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import { KubeContextService } from './kube-context.service';

describe('KubeContextService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: KubeContextService = TestBed.get(KubeContextService);
    expect(service).toBeTruthy();
  });
});
