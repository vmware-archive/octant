/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Params, Router, RouterEvent } from '@angular/router';
import { ActivatedRouteStub } from 'src/app/testing/activated-route-stub';
import { ContentComponent } from './content.component';
import { BehaviorSubject, ReplaySubject } from 'rxjs';
import { ContentResponse } from '../../../../shared/models/content';
import { ContentService } from '../../../../shared/services/content/content.service';
import { IconService } from '../../../../shared/services/icon/icon.service';
import { SugarloafModule } from '../../../sugarloaf.module';

class ContentServiceMock {
  current = new BehaviorSubject<ContentResponse>({
    content: { extensionComponent: null, viewComponents: [], title: [] },
  });
  defaultPath = new BehaviorSubject<string>('/path');
  setContentPath = (contentPath: string, params: Params) => {};
}

describe('OverviewComponent', () => {
  let component: ContentComponent;
  let fixture: ComponentFixture<ContentComponent>;
  let eventSubject: ReplaySubject<RouterEvent>;
  let routerMock;
  let contentSpy;

  beforeEach(async(() => {
    eventSubject = new ReplaySubject<RouterEvent>(1);
    routerMock = {
      events: eventSubject.asObservable(),
      routerState: {
        snapshot: {
          url: '/',
        },
      },
      parseUrl: (_: string) => {
        return {
          root: {
            children: {
              primary: {
                segments: ['foo', 'bar'],
              },
            },
          },
        };
      },
    };
    TestBed.configureTestingModule({
      imports: [SugarloafModule],
      providers: [
        { provide: ActivatedRoute, useValue: new ActivatedRouteStub({}) },
        { provide: Router, useValue: routerMock },
        { provide: ContentService, useClass: ContentServiceMock },

        IconService,
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContentComponent);

    const debugElement = fixture.debugElement;
    const contentService = debugElement.injector.get(ContentService);
    contentSpy = spyOn(contentService, 'setContentPath');

    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    const event = jasmine.createSpyObj('RoutesRecognized', ['url']);
    eventSubject.next(event);

    expect(contentSpy).toHaveBeenCalledWith('/', undefined);

    expect(component).toBeTruthy();
  });
});
