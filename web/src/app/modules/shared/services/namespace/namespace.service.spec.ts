// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { inject, TestBed } from '@angular/core/testing';
import { NamespaceService } from './namespace.service';
import {
  BackendService,
  WebsocketService,
} from '../websocket/websocket.service';
import { WebsocketServiceMock } from '../websocket/mock';

describe('NamespaceService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        NamespaceService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
      ],
    });
  });

  describe('setNamespace', () => {
    it('should tell backend service to the selected namespace', inject(
      [NamespaceService, WebsocketService],
      (svc: NamespaceService, websocketService: BackendService) => {
        spyOn(websocketService, 'sendMessage');

        svc.setNamespace('other');

        expect(websocketService.sendMessage).toHaveBeenCalledWith(
          'action.octant.dev/setNamespace',
          {
            namespace: 'other',
          }
        );
      }
    ));
  });

  describe('namespaces update', () => {
    it('triggers the list subject', inject(
      [NamespaceService, WebsocketService],
      (svc: NamespaceService, backendService: BackendService) => {
        backendService.triggerHandler('event.octant.dev/namespaces', {
          namespaces: ['foo', 'bar'],
        });
        svc.availableNamespaces.subscribe(current =>
          expect(current).toEqual(['foo', 'bar'])
        );
      }
    ));
  });
});
