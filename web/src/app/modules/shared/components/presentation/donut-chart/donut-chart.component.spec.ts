/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { DonutChartComponent } from './donut-chart.component';
import { DonutChartView } from '../../../models/content';
import { SharedModule } from '../../../shared.module';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';

describe('DonutChartComponent', () => {
  let component: DonutChartComponent;
  let fixture: ComponentFixture<DonutChartComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [DonutChartComponent, OverlayScrollbarsComponent],
        imports: [SharedModule, OverlayscrollbarsModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(DonutChartComponent);
    component = fixture.componentInstance;
    const view: DonutChartView = {
      metadata: {
        type: 'donutChart',
      },
      config: {
        segments: [
          {
            count: 3,
            status: 'ok',
          },
          {
            count: 1,
            status: 'error',
          },
        ],
        labels: {
          plural: 'items',
          singular: 'item',
        },
        size: 25,
      },
    };
    component.view = view;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should have correct segment colors', () => {
    const root: HTMLElement = fixture.nativeElement;
    const paths: NodeList = root.querySelectorAll('path');

    expect(component).toBeTruthy();
    expect(paths.length).toEqual(2);
    const el1: Element = paths[0] as Element;
    expect(el1.id).toEqual('path0');
    expect(el1.attributes.length).toEqual(5);
    expect(el1.getAttribute('fill')).toEqual('#e12200');

    const el2: Element = paths[1] as Element;
    expect(el2.id).toEqual('path1');
    expect(el2.attributes.length).toEqual(5);
    expect(el2.getAttribute('fill')).toEqual('#60b515');
  });
});
