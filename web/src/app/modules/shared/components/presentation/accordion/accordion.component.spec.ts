// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { AccordionComponent } from './accordion.component';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SharedModule } from '../../../shared.module';

describe('AccordionComponent', () => {
  let component: AccordionComponent;
  let fixture: ComponentFixture<AccordionComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [],
        imports: [SharedModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(AccordionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
