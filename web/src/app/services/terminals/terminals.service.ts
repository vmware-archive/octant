// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject } from 'rxjs';
import { TerminalOutput } from 'src/app/models/content';
import getAPIBase from '../common/getAPIBase';

const API_BASE = getAPIBase();

export class TerminalOutputStreamer {
  public scrollback: BehaviorSubject<string[]>;
  public line: BehaviorSubject<string>;

  private intervalID: number;

  constructor(
    private namespace: string,
    private pod: string,
    private container: string,
    private uuid: string,
    private http: HttpClient
  ) {}

  private poll() {
    this.http.get(this.terminalUrl()).subscribe((res: TerminalOutput) => {
      this.scrollback.next(res.scrollback);
      this.line.next(res.line);
    });
  }

  public start(): void {
    this.scrollback = new BehaviorSubject([]);
    this.line = new BehaviorSubject('');
    this.poll();
    this.intervalID = window.setInterval(() => this.poll(), 500);
  }

  public close(): void {
    this.scrollback.unsubscribe();
    clearInterval(this.intervalID);
  }

  private terminalUrl(): string {
    return [
      API_BASE,
      'api/v1',
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
  constructor(private http: HttpClient) {}

  public createStream(
    namespace,
    pod,
    container,
    uuid: string
  ): TerminalOutputStreamer {
    const pls = new TerminalOutputStreamer(
      namespace,
      pod,
      container,
      uuid,
      this.http
    );
    pls.start();
    return pls;
  }
}
