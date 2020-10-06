// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import {
  KubeContextMessage,
  KubeContextResponse,
  KubeContextService,
} from './kube-context.service';
import { WebsocketServiceMock } from '../../../../data/services/websocket/mock';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';

describe('KubeContextService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [
        KubeContextService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
      ],
    })
  );

  it('should be created', () => {
    const service: KubeContextService = TestBed.inject(KubeContextService);
    expect(service).toBeTruthy();
  });

  describe('kubeConfig update', () => {
    let service: KubeContextService;

    const update: KubeContextResponse = {
      contexts: [{ name: 'foo' }, { name: 'bar' }],
      currentContext: 'foo',
    };

    beforeEach(() => {
      service = TestBed.inject(KubeContextService);
      const backendService = TestBed.inject(WebsocketService);
      backendService.triggerHandler(KubeContextMessage, update);
    });

    it('sets the current context', () => {
      service
        .selected()
        .subscribe(selected => expect(selected).toEqual(update.currentContext));
    });

    it('sets the list of contexts', () => {
      service
        .contexts()
        .subscribe(contexts => expect(contexts).toEqual(update.contexts));
    });
  });
});
