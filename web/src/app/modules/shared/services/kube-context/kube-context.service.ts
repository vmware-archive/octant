// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';
import { take } from 'rxjs/operators';

export const KubeContextMessage = 'event.octant.dev/kubeConfig';

export interface ContextDescription {
  name: string;
}

export interface KubeContextResponse {
  contexts: ContextDescription[];
  currentContext: string;
}

const emptyKubeContext: KubeContextResponse = {
  contexts: [],
  currentContext: '',
};

@Injectable({
  providedIn: 'root',
})
export class KubeContextService {
  private contextsSource: BehaviorSubject<ContextDescription[]> =
    new BehaviorSubject<ContextDescription[]>([]);
  private selectedSource = new BehaviorSubject<string>('');

  constructor(private websocketService: WebsocketService) {
    websocketService.registerHandler(KubeContextMessage, data => {
      const update = data as KubeContextResponse;
      this.contextsSource.next(update.contexts);
      this.selectedSource.next(update.currentContext);
    });
  }

  select(context: ContextDescription) {
    this.selectedSource.pipe(take(1)).subscribe(current => {
      if (current !== context.name) {
        this.selectedSource.next(context.name);
        this.updateContext(context.name);
      }
    });
  }

  selected() {
    return this.selectedSource;
  }

  contexts() {
    return this.contextsSource;
  }

  private updateContext(name: string) {
    this.websocketService.sendMessage('action.octant.dev/setContext', {
      requestedContext: name,
    });
  }
}
