// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import _ from 'lodash';

export class EventSourceStub {
  eventListenerQueue: Array<[string, (message: MessageEvent) => void]> = [];
  eventMessageQueue: Array<[string, string]> = [];

  addEventListener(
    eventName: string,
    cb: (message: MessageEvent) => void
  ): void {
    this.eventListenerQueue.push([eventName, cb]);
  }

  close(): void {
    this.eventListenerQueue = [];
    this.eventMessageQueue = [];
  }

  queueMessage(eventName: string, data?: any): void {
    this.eventMessageQueue.push([eventName, data]);
  }

  flush(): void {
    _.remove(
      this.eventMessageQueue,
      ([messageEventName, data]): boolean => {
        const message = new MessageEvent(messageEventName, { data });
        _.forEach(this.eventListenerQueue, ([listenerEventName, cb]) => {
          if (messageEventName === listenerEventName) {
            cb(message);
          }
        });
        return true;
      }
    );
  }
}

@Injectable({
  providedIn: 'root',
})
export class EventSourceService {
  createEventSource(url: string): EventSource {
    return new EventSource(url);
  }
}
