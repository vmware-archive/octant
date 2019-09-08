// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Port } from 'src/app/models/content';
import getAPIBase from '../common/getAPIBase';
import { WebsocketService } from '../../modules/overview/services/websocket/websocket.service';

@Injectable({
  providedIn: 'root',
})
export class PortForwardService {
  constructor(private websocketService: WebsocketService) {}

  public create(port: Port) {
    const config = port.config;
    this.websocketService.sendMessage('startPortForward', {
      apiVersion: config.apiVersion,
      kind: config.kind,
      name: config.name,
      namespace: config.namespace,
      port: config.port,
    });
  }

  public remove(id: string) {
    this.websocketService.sendMessage('stopPortForward', {
      id,
    });
  }
}
