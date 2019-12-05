// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';

import { OverviewModule } from './overview.module';
import { OverviewComponent } from './overview.component';
import { BehaviorSubject } from 'rxjs';
import { ContentResponse } from '../../models/content';
import { ContentService } from './services/content/content.service';
import { IconService } from './services/icon.service';

class ContentServiceMock {
  current = new BehaviorSubject<ContentResponse>({
    content: { extensionComponent: null, viewComponents: [], title: [] },
  });
  defaultPath = new BehaviorSubject<string>('/path');
}

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
        { provide: ActivatedRoute, useValue: new ActivatedRouteStub({}) },
        { provide: Router, useValue: routerSpy },
        { provide: ContentService, useClass: ContentServiceMock },
        IconService,
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
