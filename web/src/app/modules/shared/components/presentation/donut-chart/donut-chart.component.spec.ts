/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DonutChartComponent } from './donut-chart.component';
import { DonutChartView } from '../../../models/content';
import { SharedModule } from '../../../shared.module';

describe('DonutChartComponent', () => {
  let component: DonutChartComponent;
  let fixture: ComponentFixture<DonutChartComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [DonutChartComponent],
      imports: [SharedModule],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DonutChartComponent);
    component = fixture.componentInstance;
    const view: DonutChartView = {
      metadata: {
        type: 'donutChart',
      },
      config: {
        segments: [],
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
});
