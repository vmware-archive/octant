// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { NotifierService, NotifierSession } from '../notifier/notifier.service';
import { WebsocketService } from '../../modules/overview/services/websocket/websocket.service';
import { take } from 'rxjs/operators';

interface UpdateNamespaceMessage {
  namespace: string;
}

export interface UpdateNamespacesMessage {
  namespaces: [];
}

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  activeNamespace = new BehaviorSubject<string>('');
  availableNamespaces = new BehaviorSubject<string[]>([]);

  constructor(
    private router: Router,
    private http: HttpClient,
    private websocketService: WebsocketService
  ) {
    websocketService.registerHandler('namespaces', data => {
      const update = data as UpdateNamespacesMessage;
      this.availableNamespaces.next(update.namespaces);
    });

    this.activeNamespace.subscribe(namespace => {
      websocketService.sendMessage('setNamespace', { namespace });
    });
  }

  setNamespace(namespace: string) {
    this.activeNamespace.pipe(take(1)).subscribe(cur => {
      if (cur !== namespace) {
        this.activeNamespace.next(namespace);
      }
    });
  }
}
