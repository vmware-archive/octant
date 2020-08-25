// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { TerminalComponent } from './terminal.component';
import { TerminalView } from '../../../models/content';

describe('TerminalComponent', () => {
  let component: TerminalComponent;
  let fixture: ComponentFixture<TerminalComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TerminalComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TerminalComponent);
    component = fixture.componentInstance;
    component.view = {
      metadata: {
        type: '',
      },
      config: {
        containers: [],
        name: 'name',
        namespace: 'namespace',
        podName: 'pod-name',
        terminal: {
          active: false,
        },
      },
    } as TerminalView;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
