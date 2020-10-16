/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';
import { ModalService } from './modal.service';

describe('ModalService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      declarations: [OverlayScrollbarsComponent],
      imports: [OverlayscrollbarsModule],
      providers: [ModalService],
    })
  );

  it('should be created', () => {
    const service: ModalService = TestBed.inject(ModalService);
    expect(service).toBeTruthy();
  });
});
