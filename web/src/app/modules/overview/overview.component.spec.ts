// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { ContentStreamService } from 'src/app/services/content-stream/content-stream.service';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';

import { OverviewModule } from './overview.module';
import { OverviewComponent } from './overview.component';

const contentStreamServiceStub: Partial<ContentStreamService> = {
  closeStream: () => {},
};
const activatedRouteStub = new ActivatedRouteStub();

describe('OverviewComponent', () => {
  let component: OverviewComponent;
  let fixture: ComponentFixture<OverviewComponent>;

  beforeEach(async(() => {
    const routerSpy = jasmine.createSpyObj('Router', [
      'navigateByUrl',
      'navigate',
    ]);
    TestBed.configureTestingModule({
      imports: [OverviewModule],
      providers: [
        { provide: ActivatedRoute, useValue: activatedRouteStub },
        { provide: ContentStreamService, useValue: contentStreamServiceStub },
        { provide: Router, useValue: routerSpy },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
