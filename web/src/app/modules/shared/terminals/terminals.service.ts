// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { TerminalOutput } from 'src/app/modules/shared/models/content';
import { WebsocketService } from 'src/app/modules/shared/services/websocket/websocket.service';

export class TerminalOutputStreamer {
  public line: BehaviorSubject<string>;
  public scrollback: BehaviorSubject<string>;
  public exitMessage: BehaviorSubject<string>;

  constructor(
    private namespace: string,
    private pod: string,
    private container: string,
    private wss: WebsocketService
  ) {
    this.wss.sendMessage('action.octant.dev/setActiveTerminal', {
      namespace: this.namespace,
      podName: this.pod,
      containerName: this.container,
    });

    this.line = new BehaviorSubject('');
    this.scrollback = new BehaviorSubject('');
    this.exitMessage = new BehaviorSubject('');
    this.wss.registerHandler(this.terminalUrl(), data => {
      const update = data as TerminalOutput;
      this.line.next(update.line);
      this.scrollback.next(update.scrollback);
      this.exitMessage.next(update.exitMessage);
    });
  }

  private terminalUrl(): string {
    return [
      'event.octant.dev',
      'terminals',
      `namespace/${this.namespace}`,
      `pod/${this.pod}`,
      `container/${this.container}`,
    ].join('/');
  }
}

@Injectable({
  providedIn: 'root',
})
export class TerminalOutputService {
  public selectedContainer: string;
  public namespace: string;
  public podName: string;

  constructor(private websocketService: WebsocketService) {}

  public createStream(namespace, pod, container): TerminalOutputStreamer {
    const tos = new TerminalOutputStreamer(
      namespace,
      pod,
      container,
      this.websocketService
    );
    return tos;
  }
}
