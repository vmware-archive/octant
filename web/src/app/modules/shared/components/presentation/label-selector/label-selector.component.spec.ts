// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { EditorComponent } from '../../smart/editor/editor.component';

import { LabelSelectorComponent } from './label-selector.component';

describe('LabelSelectorComponent', () => {
  let component: LabelSelectorComponent;
  let fixture: ComponentFixture<LabelSelectorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LabelSelectorComponent, EditorComponent],
      imports: [],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LabelSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
