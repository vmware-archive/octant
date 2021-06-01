// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ObjectStatusComponent } from './object-status.component';
import { SharedModule } from '../../../shared.module';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';
import { View } from '../../../models/content';

describe('ObjectStatusComponent', () => {
  let component: ObjectStatusComponent;
  let fixture: ComponentFixture<ObjectStatusComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [OverlayScrollbarsComponent],
        imports: [SharedModule, OverlayscrollbarsModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ObjectStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show properties table', () => {
    component.node = {
      apiVersion: 'v1',
      details: [],
      kind: 'Pod',
      name: 'test-pod',
      path: undefined,
      properties: [
        {
          label: 'test',
          value: {
            metadata: { type: 'text' },
            config: { value: 'property text' },
          } as View,
        },
      ],
      status: 'pod status',
    };
    fixture.detectChanges();

    const root: HTMLElement = fixture.nativeElement;
    const el: SVGPathElement = root.querySelector('.properties');

    expect(component).toBeTruthy();
    expect(el).not.toBeNull();
    expect(el.innerHTML).toContain('Pod');
    expect(el.innerHTML).toContain('v1');
    expect(el.innerHTML).toContain('property text');
  });
});
