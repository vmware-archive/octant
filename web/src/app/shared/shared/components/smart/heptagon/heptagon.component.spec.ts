// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { HeptagonComponent } from './heptagon.component';
import { Point } from '../../../../../modules/overview/models/point';

describe('HeptagonComponent', () => {
  let component: HeptagonComponent;
  let fixture: ComponentFixture<HeptagonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [HeptagonComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(HeptagonComponent);
    component = fixture.componentInstance;

    component.status = {
      status: 'OK',
      name: 'name',
    };
    component.centerPoint = new Point(10, 10);
    component.edgeLength = 7;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have a path element with the heptagon points', () => {
    const root: HTMLElement = fixture.nativeElement;
    const el: SVGPathElement = root.querySelector('path');

    const expectedPoints = [
      'M17.864428613011135,11.795004510726969',
      'L13.499999999999998,17.26782488800318',
      'L6.5,17.26782488800318',
      'L2.1355713869888646,11.79500451072697',
      'L3.6932179246830676,4.9705091254542015',
      'L9.999999999999998,1.9333229516312969',
      'L16.306782075316928,4.970509125454197',
    ];

    expect(el.getAttribute('d')).toEqual(expectedPoints.join(' '));
  });

  it('should set the heptagon class', () => {
    const root: HTMLElement = fixture.nativeElement;
    const el: SVGPathElement = root.querySelector('path');

    expect(el.classList.contains('heptagon')).toBeTruthy();
  });

  it('should set the pod status', () => {
    const root: HTMLElement = fixture.nativeElement;
    const el: SVGPathElement = root.querySelector('path');

    expect(el.classList.contains('status-OK')).toBeTruthy();
  });
});
