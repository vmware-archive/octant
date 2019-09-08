/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { inject, TestBed } from '@angular/core/testing';

import { BackendService, WebsocketService } from './websocket.service';
import {
  NotifierService,
  NotifierSession,
  NotifierSignal,
} from '../../../../services/notifier/notifier.service';
import _ from 'lodash';
import { BehaviorSubject } from 'rxjs';

class NotifierServiceMock {
  private signalsStream: BehaviorSubject<NotifierSignal[]>;

  createSession = (): NotifierSession => {
    return new NotifierSession(this.signalsStream, _.uniqueId('signalSession'));
  };
}

describe('WebsocketService', () => {
  beforeEach(() =>
    TestBed.configureTestingModule({
      providers: [
        WebsocketService,
        {
          provide: NotifierService,
          useClass: NotifierServiceMock,
        },
      ],
    })
  );

  describe('sendMessage', () => {
    it('should send a message', inject([WebsocketService], service => {
      service.sendMessage('messageType', {});
    }));
  });

  it('should be created', () => {
    const service: WebsocketService = TestBed.get(WebsocketService);
    expect(service).toBeTruthy();
  });
});
