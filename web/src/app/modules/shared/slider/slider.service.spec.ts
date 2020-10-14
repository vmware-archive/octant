// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed, waitForAsync } from '@angular/core/testing';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';
import { EditorComponent } from '../components/smart/editor/editor.component';
import { SliderService } from './slider.service';

describe('SliderService', () => {
  let service: SliderService;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [OverlayScrollbarsComponent, EditorComponent],
        providers: [OverlayscrollbarsModule],
      });
      service = TestBed.inject(SliderService);
    })
  );

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it(
    'set value',
    waitForAsync(() => {
      service.setHeight(100);
      service.setHeight$.subscribe(current => expect(current).toEqual(100));
    })
  );

  it(
    'reset to default',
    waitForAsync(() => {
      service.resetDefault();
      service.setHeight$.subscribe(current => expect(current).toEqual(36));
    })
  );
});
