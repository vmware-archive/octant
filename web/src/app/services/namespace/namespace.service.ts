// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import { NotifierService, NotifierSession } from '../notifier/notifier.service';
import { WebsocketService } from '../../modules/overview/services/websocket/websocket.service';

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
  activeNamespace = new BehaviorSubject<string>('default');
  availableNamespaces = new BehaviorSubject<string[]>([]);

  constructor(
    private router: Router,
    private http: HttpClient,
    private websocketService: WebsocketService
  ) {
    websocketService.registerHandler('namespace', data => {
      const update = data as UpdateNamespaceMessage;
      this.activeNamespace.next(update.namespace);
    });

    websocketService.registerHandler('namespaces', data => {
      const update = data as UpdateNamespacesMessage;
      this.availableNamespaces.next(update.namespaces);
    });

    this.activeNamespace.subscribe(namespace => {
      websocketService.sendMessage('setNamespace', { namespace });
    });
  }

  setNamespace(namespace: string) {
    this.activeNamespace.next(namespace);
  }
}
