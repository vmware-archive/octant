/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Inject, Injectable } from '@angular/core';
import { Router } from '@angular/router';
import {
  BehaviorSubject,
  ObjectUnsubscribedError,
  Observable,
  Subject,
} from 'rxjs';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../../../modules/shared/notifier/notifier.service';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { delay, retryWhen, take, tap } from 'rxjs/operators';
import { WindowToken } from '../../../window';
import { ElectronService } from 'src/app/modules/shared/services/electron/electron.service';

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
  handlers: { [key: string]: ({}) => void } = {};
  reconnected = new Subject<Event>();

  private notifierSession: NotifierSession;
  private subject: WebSocketSubject<unknown>;
  private connectSignalID = new BehaviorSubject<string>('');
  private opened = false;

  private router: Router;

  constructor(
    private electronService: ElectronService,
    notifierService: NotifierService,
    router: Router,
    @Inject(WindowToken) private window: Window
  ) {
    this.notifierSession = notifierService.createSession();
    this.router = router;

    this.registerHandler('event.octant.dev/alert', data => {
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

    this.registerHandler('event.octant.dev/contentPath', data => {
      const payload = data as { contentPath: string };
      const contentPath = payload.contentPath || '/';
      if (this.router.url !== contentPath) {
        this.router.navigateByUrl(contentPath);
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
    if (this.opened) {
      return;
    }
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
          this.connectSignalID.pipe(take(1)).subscribe(id => {
            if (id !== '') {
              this.notifierSession.removeAllSignals();
              this.connectSignalID.next('');
            }
          });

          this.parseWebsocketMessage(data);
        },
        err => console.error(err),
        () => {
          console.log('web socket is closing');
        }
      );

    this.opened = true;
  }

  close() {
    if (!this.opened) {
      this.opened = false;
    }
    this.subject.unsubscribe();
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
      if (!(err instanceof ObjectUnsubscribedError)) {
        console.error('parse websocket', err, data);
      }
    }
  }

  websocketURI(): string {
    // TODO: https://github.com/vmware-tanzu/octant/issues/944
    if (this.electronService.isElectron()) {
      return 'ws://localhost:7777/api/v1/stream';
    }

    const loc = this.window.location;
    let newURI = 'ws:';
    if (loc.protocol === 'https:') {
      newURI = 'wss:';
    }
    newURI += '//' + loc.host;
    newURI += loc.pathname + 'api/v1/stream';
    return newURI;
  }
}
