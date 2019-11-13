// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { TerminalOutput } from 'src/app/models/content';
import { WebsocketService } from 'src/app/modules/overview/services/websocket/websocket.service';

export class TerminalOutputStreamer {
  public line: BehaviorSubject<string>;
  public scrollback: BehaviorSubject<string>;

  constructor(
    private namespace: string,
    private pod: string,
    private container: string,
    private uuid: string,
    private wss: WebsocketService
  ) {
    this.wss.sendMessage('sendTerminalScrollback', {
      terminalID: this.uuid,
    });

    this.line = new BehaviorSubject('');
    this.scrollback = new BehaviorSubject('');
    this.wss.registerHandler(this.terminalUrl(), data => {
      const update = data as TerminalOutput;
      this.line.next(update.line);
      this.scrollback.next(update.scrollback);
    });
  }

  private terminalUrl(): string {
    return [
      'terminals',
      `namespace/${this.namespace}`,
      `pod/${this.pod}`,
      `container/${this.container}`,
      `${this.uuid}`,
    ].join('/');
  }
}

@Injectable({
  providedIn: 'root',
})
export class TerminalOutputService {
  constructor(private websocketService: WebsocketService) {}

  public createStream(
    namespace,
    pod,
    container,
    uuid: string
  ): TerminalOutputStreamer {
    const tos = new TerminalOutputStreamer(
      namespace,
      pod,
      container,
      uuid,
      this.websocketService
    );
    return tos;
  }
}
