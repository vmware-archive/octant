// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ResourceViewerComponent } from './resource-viewer.component';
import { SharedModule } from '../../../shared.module';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';

describe('ResourceViewerComponent', () => {
  let component: ResourceViewerComponent;
  let fixture: ComponentFixture<ResourceViewerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [OverlayScrollbarsComponent],
      imports: [SharedModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceViewerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
