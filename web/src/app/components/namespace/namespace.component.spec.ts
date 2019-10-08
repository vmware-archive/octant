// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { NamespaceComponent } from './namespace.component';
import { NgSelectModule } from '@ng-select/ng-select';

describe('NamespaceComponent', () => {
  let component: NamespaceComponent;
  let fixture: ComponentFixture<NamespaceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [NgSelectModule],
      declarations: [NamespaceComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NamespaceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
