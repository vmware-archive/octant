// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { HeptagonGridComponent } from './heptagon-grid.component';
import {
  HeptagonGridRowComponent,
  HoverStatus,
} from '../heptagon-grid-row/heptagon-grid-row.component';
import { HeptagonLabelComponent } from '../heptagon-label/heptagon-label.component';
import { HeptagonComponent } from '../heptagon/heptagon.component';

describe('HeptagonGridComponent', () => {
  let component: HeptagonGridComponent;
  let fixture: ComponentFixture<HeptagonGridComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        HeptagonGridComponent,
        HeptagonGridRowComponent,
        HeptagonLabelComponent,
        HeptagonComponent,
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HeptagonGridComponent);
    component = fixture.componentInstance;

    component.podStatuses = [
      { name: 'pod-1', status: 'ok' },
      { name: 'pod-2', status: 'ok' },
      { name: 'pod-3', status: 'ok' },
      { name: 'pod-4', status: 'ok' },
      { name: 'pod-5', status: 'ok' },
    ];
    component.perRow = 2;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('updates hover states', () => {
    const status: HoverStatus = {
      col: 1,
      row: 1,
      hovered: true,
    };
    component.updateHover(status);

    expect(component.hoverStates[1][1]).toBeTruthy();
  });

  describe('heptagon is hovered', () => {
    beforeEach(() => {
      component.hoverStates[0][0] = true;
    });

    it('knows if a heptagon is activated', () => {
      expect(component.isActivated(0)).toBeTruthy();
    });

    it('know if a heptagon is not activated', () => {
      expect(component.isActivated(1)).toBeFalsy();
    });
  });
});
