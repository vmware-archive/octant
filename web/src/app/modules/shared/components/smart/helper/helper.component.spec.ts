// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HelperComponent } from './helper.component';
import { HelperService } from '../../../services/helper/helper.service';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

describe('HelperComponent', () => {
  let component: HelperComponent;
  let fixture: ComponentFixture<HelperComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [HelperComponent],
      providers: [{ provide: HelperService }],
      imports: [BrowserAnimationsModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HelperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
