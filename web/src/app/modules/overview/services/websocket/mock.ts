/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { BackendService, HandlerFunc } from './websocket.service';

export class WebsocketServiceMock implements BackendService {
  private handlers: { [key: string]: HandlerFunc } = {};

  isOpen = false;

  sendMessage = (messageType: string, payload: {}) => {};

  close() {
    this.isOpen = false;
  }

  open() {
    this.isOpen = true;
  }

  registerHandler(name: string, handler: HandlerFunc) {
    this.handlers[name] = handler;
  }

  triggerHandler(name: string, payload: {}) {
    if (!this.handlers[name]) {
      throw new Error(`handler ${name} was not found`);
    }
    this.handlers[name](payload);
  }
}
