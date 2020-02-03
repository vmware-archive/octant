// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContentSwitcherComponent } from './content-switcher.component';
import { SharedModule } from '../../../shared.module';

describe('ContentSwitcherComponent', () => {
  let component: ContentSwitcherComponent;
  let fixture: ComponentFixture<ContentSwitcherComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [SharedModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContentSwitcherComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
