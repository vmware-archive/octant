// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Router } from '@angular/router';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';

export interface BuildInfoMessage {
  version: string;
  commit: string;
  time: string;
}

@Injectable({
  providedIn: 'root',
})
export class HelperService {
  private version = new BehaviorSubject<string>('');
  private commit = new BehaviorSubject<string>('');
  private time = new BehaviorSubject<string>('');

  constructor(
    private router: Router,
    private websocketService: WebsocketService
  ) {
    websocketService.registerHandler('event.octant.dev/buildInfo', data => {
      const update = data as BuildInfoMessage;
      this.version.next(update.version);
      this.commit.next(update.commit);
      this.time.next(update.time);
    });
  }

  buildVersion() {
    return this.version;
  }

  buildCommit() {
    return this.commit;
  }

  buildTime() {
    return this.time;
  }
}
