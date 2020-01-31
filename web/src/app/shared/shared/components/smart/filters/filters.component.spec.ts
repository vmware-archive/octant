// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import {
  Filter,
  LabelFilterService,
} from 'src/app/services/label-filter/label-filter.service';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';
import { FormsModule } from '@angular/forms';
import { OverviewModule } from '../../../../../modules/overview/overview.module';
import { FiltersComponent } from './filters.component';

const filterSubject = new BehaviorSubject<Filter[]>([]);
const labelFilterService: Partial<LabelFilterService> = {
  filters: filterSubject,
};

const activatedRouteStub = new ActivatedRouteStub();

describe('FiltersComponent', () => {
  let component: FiltersComponent;
  let fixture: ComponentFixture<FiltersComponent>;
  let routerSpy: any;

  beforeEach(async(() => {
    const mockRouter = {
      navigate: jasmine.createSpy('navigate'),
    };

    TestBed.configureTestingModule({
      imports: [OverviewModule, FormsModule],
      providers: [
        { provide: Router, useValue: mockRouter },
        { provide: ActivatedRoute, useValue: activatedRouteStub },
        { provide: LabelFilterService, useValue: labelFilterService },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    routerSpy = TestBed.get(Router);
    fixture = TestBed.createComponent(FiltersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', async(() => {
    fixture.whenStable().then(() => {
      expect(component).toBeTruthy();
    });
  }));
});
