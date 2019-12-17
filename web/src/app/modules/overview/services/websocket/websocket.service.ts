/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, Subject, Subscription } from 'rxjs';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../../../../services/notifier/notifier.service';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { delay, retryWhen, tap } from 'rxjs/operators';

interface WebsocketPayload {
  type: string;
  data?: {};
}

export type HandlerFunc = (data: {}) => void;

export interface BackendService {
  open();
  close();
  registerHandler(name: string, handler: HandlerFunc);
  sendMessage(messageType: string, payload: {});
  triggerHandler(name: string, payload: {});
}

interface Alert {
  type: NotifierSignalType;
  message: string;
  expiration?: string;
}

@Injectable({
  providedIn: 'root',
})
export class WebsocketService implements BackendService {
  ws: WebSocket;
  handlers: { [key: string]: ({}) => void } = {};
  reconnected = new Subject<Event>();

  private notifierSession: NotifierSession;
  private subject: WebSocketSubject<unknown>;
  private connectSignalID = new BehaviorSubject<string>('');

  constructor(notifierService: NotifierService) {
    this.notifierSession = notifierService.createSession();

    this.registerHandler('alert', data => {
      const alert = data as Alert;
      const id = this.notifierSession.pushSignal(alert.type, alert.message);
      if (alert.expiration) {
        const expiration = new Date(alert.expiration);
        const diff = expiration.getTime() - Date.now();

        setTimeout(() => {
          this.notifierSession.removeSignal(id);
        }, diff);
      }
    });
  }

  registerHandler(name: string, handler: (data: {}) => void): () => void {
    this.handlers[name] = handler;
    return () => delete this.handlers[name];
  }

  triggerHandler(name: string, payload: {}) {
    if (!this.handlers[name]) {
      throw new Error(`handler ${name} was not found`);
    }
    this.handlers[name](payload);
  }

  open() {
    this.createWebSocket()
      .pipe(
        retryWhen(errors =>
          errors.pipe(
            tap(_ => {
              const id = this.notifierSession.pushSignal(
                NotifierSignalType.ERROR,
                'Lost connection to Octant service. Retrying...'
              );
              this.connectSignalID.next(id);
            }),
            delay(1000)
          )
        )
      )
      .subscribe(
        data => {
          this.connectSignalID
            .subscribe(id => {
              if (id !== '') {
                this.notifierSession.removeAllSignals();
                this.connectSignalID.next('');
              }
            })
            .unsubscribe();

          this.parseWebsocketMessage(data);
        },
        err => console.error(err)
      );
  }

  close() {
    this.subject.unsubscribe();
  }

  private websocketURI() {
    const loc = window.location;
    let newURI = '';
    if (loc.protocol === 'https:') {
      newURI = 'wss:';
    } else {
      newURI = 'ws:';
    }
    newURI += '//' + loc.host;
    newURI += loc.pathname + 'api/v1/stream';
    return newURI;
  }

  private createWebSocket() {
    const uri = this.websocketURI();
    return new Observable(observer => {
      try {
        const subject = webSocket({
          url: uri,
          deserializer: ({ data }) => JSON.parse(data),
          openObserver: this.reconnected,
        });

        const subscription = subject.asObservable().subscribe(
          data => observer.next(data),
          error => observer.error(error),
          () => observer.complete()
        );

        this.subject = subject;
        return () => {
          if (!subscription.closed) {
            subscription.unsubscribe();
          }
        };
      } catch (error) {
        observer.error(error);
      }
    });
  }

  sendMessage(messageType: string, payload: {}) {
    if (this.subject) {
      const data = {
        type: messageType,
        payload,
      };
      this.subject.next(data);
    }
  }

  private parseWebsocketMessage(data: {}) {
    try {
      const payload = data as WebsocketPayload;
      if (this.handlers.hasOwnProperty(payload.type)) {
        const handler = this.handlers[payload.type];
        handler(payload.data);
      } else {
        console.warn(
          `received websocket unknown message of type ${payload.type} with`,
          payload.data
        );
      }
    } catch (err) {
      console.error('parse websocket', err, data);
    }
  }
}
