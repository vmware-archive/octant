// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HeptagonGridComponent } from '../heptagon-grid/heptagon-grid.component';
import { PodStatusComponent } from './pod-status.component';
import { Component, Input } from '@angular/core';
import { PodStatus } from '../../../models/pod-status';
import { HighlightModule, HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';

@Component({
  selector: 'app-heptagon-grid',
  template: ``,
})
class TestGridComponent {
  @Input()
  podStatuses: PodStatus[] = [];

  @Input()
  edgeLength: number;
}

describe('PodStatusComponent', () => {
  let component: PodStatusComponent;
  let fixture: ComponentFixture<PodStatusComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [HighlightModule],
      declarations: [PodStatusComponent, TestGridComponent],
      providers: [
        { provide: HeptagonGridComponent, useClass: TestGridComponent },
        {
          provide: HIGHLIGHT_OPTIONS,
          useValue: {
            languages: {
              json: () => import('highlight.js/lib/languages/json'),
              yaml: () => import('highlight.js/lib/languages/yaml'),
            },
          },
        },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PodStatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
