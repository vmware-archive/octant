/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';

import { WebsocketService } from './websocket.service';
import {
  NotifierService,
  NotifierSession,
  NotifierSignal,
} from '../../../modules/shared/notifier/notifier.service';
import uniqueId from 'lodash/uniqueId';
import { BehaviorSubject } from 'rxjs';
import { WindowToken } from '../../../window';

class NotifierServiceMock {
  private signalsStream: BehaviorSubject<NotifierSignal[]>;

  createSession = (): NotifierSession => {
    return new NotifierSession(this.signalsStream, uniqueId('signalSession'));
  };
}

interface Location {
  protocol: string;
  host: string;
  pathname: string;
}

describe('WebsocketService', () => {
  const location: Location = {
    protocol: 'http',
    host: 'example.com',
    pathname: '/path/',
  };
  let service: WebsocketService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        WebsocketService,
        {
          provide: NotifierService,
          useClass: NotifierServiceMock,
        },
        {
          provide: WindowToken,
          useValue: {
            location,
          },
        },
      ],
    });

    service = TestBed.inject(WebsocketService);
  });

  describe('websocketURI', () => {
    describe('with http location', () => {
      beforeEach(() => {
        location.protocol = 'http:';
      });

      it('returns a http websocket uri', () => {
        expect(service.websocketURI()).toEqual(
          'ws://example.com/path/api/v1/stream'
        );
      });
    });
    describe('with https location', () => {
      beforeEach(() => {
        location.protocol = 'https:';
      });

      it('returns a https websocket uri', () => {
        expect(service.websocketURI()).toEqual(
          'wss://example.com/path/api/v1/stream'
        );
      });
    });
  });
});
