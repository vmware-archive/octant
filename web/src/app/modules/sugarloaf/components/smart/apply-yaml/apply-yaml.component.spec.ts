// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ApplyYAMLComponent } from './apply-yaml.component';

describe('ApplyYAMLComponent', () => {
  let component: ApplyYAMLComponent;
  let fixture: ComponentFixture<ApplyYAMLComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ApplyYAMLComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ApplyYAMLComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
