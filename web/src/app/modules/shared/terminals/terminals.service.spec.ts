// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';
import { TerminalOutputService } from './terminals.service';
import { WebsocketService } from 'src/app/data/services/websocket/websocket.service';
import { WebsocketServiceMock } from 'src/app/data/services/websocket/mock';

describe('TerminalOutputService', () => {
  let service: TerminalOutputService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        TerminalOutputService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
      ],
    });

    service = TestBed.inject(TerminalOutputService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
