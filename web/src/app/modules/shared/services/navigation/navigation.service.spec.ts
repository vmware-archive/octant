/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { inject, TestBed } from '@angular/core/testing';

import { NavigationService } from './navigation.service';
import {
  BackendService,
  WebsocketService,
} from '../../../../data/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../../data/services/websocket/mock';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';
import {
  expectedSelection,
  NAVIGATION_MOCK_DATA,
} from './navigation.test.data';
import { BehaviorSubject } from 'rxjs';
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
              if (child.children) {
                child.children.map(grandchild => {
                  verifySelection(grandchild.path, svc, index, 'grandchild');
                });
              }
            });
          }
        });
        svc.activeUrl.unsubscribe();
        svc.selectedItem.unsubscribe();
      }
    ));

    function verifySelection(
      path: string,
      svc: NavigationService,
      index: number,
      descriptor: string
    ) {
      const prefixedPath = path.startsWith('/') ? path : '/' + path;
      svc.activeUrl.next(prefixedPath);
      svc.updateLastSelection();

      svc.selectedItem.pipe(take(1)).subscribe(selection => {
        const expected = expectedSelection[path];

        expect(selection.index)
          .withContext(
            `navigation selected ${descriptor} index ${index} ${path}`
          )
          .toEqual(expected.index);
        expect(selection.module)
          .withContext(
            `navigation selected ${descriptor} module ${index} ${path}`
          )
          .toEqual(expected.module);
      });

      svc.activeUrl
        .pipe(take(1))
        .subscribe(url =>
          expect(url)
            .withContext(`url path for ${descriptor} index ${index}`)
            .toEqual(prefixedPath)
        );
    }

    const routerLinkCases = [
      { url: '/', namespace: 'test', result: '/' },
      {
        url: '/workloads/namespace/default',
        namespace: 'test',
        result: '/workloads/namespace/test',
      },
      {
        url: '/cluster-overview',
        namespace: 'test',
        result: '/cluster-overview',
      },
      { url: '/plugin/path', namespace: 'test', result: '/plugin/path' },
      {
        url: '/overview/namespace/default',
        namespace: 'test',
        result: '/overview/namespace/test',
      },
      {
        url: '/overview/namespace/default/workloads/deployments',
        namespace: 'test',
        result: '/overview/namespace/test/workloads/deployments',
      },
      {
        url:
          '/overview/namespace/default/workloads/deployments/nginx-deployment',
        namespace: 'test',
        result: '/overview/namespace/test',
      },
    ];

    routerLinkCases.forEach((test, index) => {
      it(`generates correct routerLink based on url ${test.url}`, inject(
        [NavigationService],
        (svc: NavigationService) => {
          svc.activeUrl = new BehaviorSubject<string>(test.url);
          const result = svc.redirect(test.namespace);
          expect(test.result).toEqual(result);
        }
      ));
    });
  });
});
