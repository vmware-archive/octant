// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DefaultPipe } from '../../../../shared/pipes/default/default.pipe';
import { QuickSwitcherComponent } from './quick-switcher.component';
import { windowProvider, WindowToken } from '../../../../../window';

describe('QuickSwitcherComponent', () => {
  let component: QuickSwitcherComponent;
  let fixture: ComponentFixture<QuickSwitcherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [QuickSwitcherComponent, DefaultPipe],
      providers: [{ provide: WindowToken, useFactory: windowProvider }],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(QuickSwitcherComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
