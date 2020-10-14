// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import {
  Filter,
  LabelFilterService,
} from 'src/app/modules/shared/services/label-filter/label-filter.service';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';
import { FiltersComponent } from './filters.component';
import { SharedModule } from '../../../shared.module';

const filterSubject = new BehaviorSubject<Filter[]>([]);
const labelFilterService: Partial<LabelFilterService> = {
  filters: filterSubject,
};

const activatedRouteStub = new ActivatedRouteStub();

describe('FiltersComponent', () => {
  let component: FiltersComponent;
  let fixture: ComponentFixture<FiltersComponent>;
  let routerSpy: any;

  beforeEach(
    waitForAsync(() => {
      const mockRouter = {
        navigate: jasmine.createSpy('navigate'),
      };

      TestBed.configureTestingModule({
        declarations: [ApplyYAMLComponent, OverlayScrollbarsComponent],
        imports: [SharedModule],
        providers: [
          { provide: Router, useValue: mockRouter },
          { provide: ActivatedRoute, useValue: activatedRouteStub },
          { provide: LabelFilterService, useValue: labelFilterService },
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    routerSpy = TestBed.inject(Router);
    fixture = TestBed.createComponent(FiltersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it(
    'should create',
    waitForAsync(() => {
      fixture.whenStable().then(() => {
        expect(component).toBeTruthy();
      });
    })
  );
});
