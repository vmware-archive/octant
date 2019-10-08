// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { Point } from '../../models/point';
import { HeptagonLabelComponent } from './heptagon-label.component';

describe('HeptagonLabelComponent', () => {
  let component: HeptagonLabelComponent;
  let fixture: ComponentFixture<HeptagonLabelComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [HeptagonLabelComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HeptagonLabelComponent);
    component = fixture.componentInstance;
    component.centerPoint = new Point(20, 20);
    component.height = 20;
    component.status = {
      status: 'ok',
      name: 'pod-name',
    };
    component.name = 'name';
    fixture.detectChanges();
  });

  it('sets the label container dimensions', () => {
    const el = component.container.nativeElement;
    expect(el.getAttribute('x')).toEqual('-25');
    expect(el.getAttribute('y')).toEqual('5');
    expect(el.getAttribute('width')).toEqual('30px');
    expect(el.getAttribute('height')).toEqual('20');
  });

  it('knows the font size', () => {
    expect(component.fontSize()).toEqual(14);
  });
});
