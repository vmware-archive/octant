// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../../notifier/notifier.service';
import { WebsocketService } from '../websocket/websocket.service';
import { take } from 'rxjs/operators';

export interface UpdateNamespacesMessage {
  namespaces: [];
}

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  private notifierSession: NotifierSession;

  activeNamespace = new BehaviorSubject<string>('');
  availableNamespaces = new BehaviorSubject<string[]>([]);

  constructor(
    private router: Router,
    private http: HttpClient,
    private websocketService: WebsocketService,
    private notifierService: NotifierService
  ) {
    websocketService.registerHandler('event.octant.dev/namespaces', data => {
      const update = data as UpdateNamespacesMessage;
      this.availableNamespaces.next(update.namespaces);

      this.validateNamespace(update.namespaces);
    });

    this.activeNamespace.subscribe(namespace => {
      if (namespace.length > 0) {
        websocketService.sendMessage('action.octant.dev/setNamespace', {
          namespace,
        });
      }
    });

    this.notifierSession = this.notifierService.createSession();
  }

  setNamespace(namespace: string) {
    this.activeNamespace.pipe(take(1)).subscribe(cur => {
      if (cur !== namespace) {
        this.activeNamespace.next(namespace);
        this.validateNamespace(this.availableNamespaces.getValue());
      }
    });
  }

  validateNamespace(namespaces: string[]) {
    this.activeNamespace.pipe(take(1)).subscribe(cur => {
      if (!namespaces.includes(cur) && cur !== '') {
        this.notifierSession.pushSignal(
          NotifierSignalType.WARNING,
          'Namespace does not exist.'
        );
      } else {
        this.notifierSession.removeAllSignals();
      }
    });
  }
}
