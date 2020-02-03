/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { inject, TestBed } from '@angular/core/testing';

import { NavigationService } from './navigation.service';
import {
  BackendService,
  WebsocketService,
} from '../websocket/websocket.service';
import { WebsocketServiceMock } from '../websocket/mock';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';

describe('NavigationService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [
        NavigationService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
        ContentService,
      ],
    })
  );

  it('should be created', () => {
    const service: NavigationService = TestBed.get(NavigationService);
    expect(service).toBeTruthy();
  });

  describe('namespace update', () => {
    it('triggers the current subject', inject(
      [NavigationService, WebsocketService, ContentService],
      (
        svc: NavigationService,
        backendService: BackendService,
        contentService: ContentService
      ) => {
        const update: Navigation = {
          sections: [],
          defaultPath: 'path',
        };

        backendService.triggerHandler('navigation', update);
        svc.current.subscribe(current => expect(current).toEqual(update));
        contentService.defaultPath.subscribe(defaultPath =>
          expect(defaultPath).toEqual(update.defaultPath)
        );
      }
    ));
  });
});
