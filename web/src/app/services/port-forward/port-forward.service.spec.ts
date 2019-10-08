// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';

import { PortForwardService } from './port-forward.service';
import { WebsocketService } from '../../modules/overview/services/websocket/websocket.service';
import { Port } from '../../models/content';

describe('PortForwardService', () => {
  let service: PortForwardService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PortForwardService, WebsocketService],
    });

    service = TestBed.get(PortForwardService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('create port forward', () => {
    let websocketService: WebsocketService;

    const port: Port = {
      config: {
        apiVersion: 'apiVersion',
        kind: 'kind',
        name: 'name',
        namespace: 'namespace',
        port: 1234,
        state: undefined,
        protocol: '', // no support for protocol
      },
      metadata: undefined,
    };

    beforeEach(() => {
      websocketService = TestBed.get(WebsocketService);
      spyOn(websocketService, 'sendMessage');

      service.create(port);
    });

    it('sends a websocket port forward start message', () => {
      expect(websocketService.sendMessage).toHaveBeenCalledWith(
        'startPortForward',
        {
          apiVersion: port.config.apiVersion,
          kind: port.config.kind,
          name: port.config.name,
          namespace: port.config.namespace,
          port: port.config.port,
        }
      );
    });
  });

  describe('remove port forward', () => {
    let websocketService: WebsocketService;

    const id = 'id';

    beforeEach(() => {
      websocketService = TestBed.get(WebsocketService);
      spyOn(websocketService, 'sendMessage');

      service.remove(id);
    });

    it('sends a websocket port forward stop message', () => {
      expect(websocketService.sendMessage).toHaveBeenCalledWith(
        'stopPortForward',
        {
          id,
        }
      );
    });
  });
});
