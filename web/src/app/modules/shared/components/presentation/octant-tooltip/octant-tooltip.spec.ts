// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { OctantTooltipComponent } from './octant-tooltip';
import { SharedModule } from '../../../shared.module';
import { windowProvider, WindowToken } from '../../../../../window';

describe('OctantTooltipComponent', () => {
  let component: OctantTooltipComponent;
  let fixture: ComponentFixture<OctantTooltipComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [SharedModule],
      providers: [{ provide: WindowToken, useFactory: windowProvider }],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OctantTooltipComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
