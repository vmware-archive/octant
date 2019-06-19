// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { GraphvizComponent } from './graphviz.component';

describe('GraphvizComponent', () => {
  let component: GraphvizComponent;
  let fixture: ComponentFixture<GraphvizComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ GraphvizComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(GraphvizComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
