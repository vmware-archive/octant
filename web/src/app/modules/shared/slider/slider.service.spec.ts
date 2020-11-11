// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, TestBed } from '@angular/core/testing';
import { EditorComponent } from '../components/smart/editor/editor.component';
import { SliderService } from './slider.service';

describe('SliderService', () => {
  let service: SliderService;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [EditorComponent],
      providers: [],
    });
    service = TestBed.inject(SliderService);
  }));

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('set value', async(() => {
    service.setHeight(100);
    service.setHeight$.subscribe(current => expect(current).toEqual(100));
  }));

  it('reset to default', async(() => {
    service.resetDefault();
    service.setHeight$.subscribe(current => expect(current).toEqual(36));
  }));
});
