// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { SelectorsComponent } from './selectors.component';
import { SharedModule } from '../../../shared.module';
import { EditorComponent } from '../../smart/editor/editor.component';

describe('SelectorsComponent', () => {
  let component: SelectorsComponent;
  let fixture: ComponentFixture<SelectorsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [EditorComponent],
      imports: [SharedModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectorsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
