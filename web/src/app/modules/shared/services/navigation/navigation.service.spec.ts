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
import { NAVIGATION_MOCK_DATA } from './navigation.test.data';
import { take } from 'rxjs/operators';

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
    const service: NavigationService = TestBed.inject(NavigationService);
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

        backendService.triggerHandler('event.octant.dev/navigation', update);
        svc.current.subscribe(current => expect(current).toEqual(update));
      }
    ));

    it('verify nav selection is correct after going to new URL', inject(
      [NavigationService, WebsocketService, ContentService],
      (svc: NavigationService, backendService: BackendService) => {
        const currentNavigation: Navigation = {
          sections: NAVIGATION_MOCK_DATA,
          defaultPath: '',
        };

        backendService.triggerHandler(
          'event.octant.dev/navigation',
          currentNavigation
        );

        currentNavigation.sections.map((section, index) => {
          verifySelection(section.path, svc, index, '');

          if (section.children) {
            section.children.map(child => {
              verifySelection(child.path, svc, index, 'child');
            });
          }
        });
        svc.activeUrl.unsubscribe();
        svc.lastSelection.unsubscribe();
      }
    ));

    function verifySelection(
      path: string,
      svc: NavigationService,
      index: number,
      descriptor: string
    ) {
      svc.activeUrl.next('/' + path);
      svc.updateLastSelection();

      svc.lastSelection.pipe(take(1)).subscribe(selection => {
        expect(selection)
          .withContext(`navigation selected ${descriptor} index ${index}`)
          .toEqual(index);
      });

      svc.activeUrl.pipe(take(1)).subscribe(url =>
        expect(url)
          .withContext(`url path for ${descriptor} index ${index}`)
          .toEqual('/' + path)
      );
    }
  });
});
