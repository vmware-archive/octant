// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject } from 'rxjs';
import { LogEntry, LogResponse } from 'src/app/models/content';
import getAPIBase from '../common/getAPIBase';

const API_BASE = getAPIBase();

export class PodLogsStreamer {
  public logEntries: BehaviorSubject<LogEntry[]>;
  private intervalID: number;

  constructor(
    private namespace: string,
    private pod: string,
    private container: string,
    private http: HttpClient
  ) {}

  private poll() {
    this.http.get(this.logsUrl()).subscribe((res: LogResponse) => {
      this.logEntries.next(res.entries);
    });
  }

  public start(): void {
    this.logEntries = new BehaviorSubject([]);
    this.poll();
    this.intervalID = window.setInterval(() => this.poll(), 5000);
  }

  public close(): void {
    this.logEntries.unsubscribe();
    clearInterval(this.intervalID);
  }

  private logsUrl(): string {
    return [
      API_BASE,
      'api/v1',
      'logs',
      `namespace/${this.namespace}`,
      `pod/${this.pod}`,
      `container/${this.container}`,
    ].join('/');
  }
}

@Injectable({
  providedIn: 'root',
})
export class PodLogsService {
  constructor(private http: HttpClient) {}

  public createStream(namespace, pod, container: string): PodLogsStreamer {
    const pls = new PodLogsStreamer(namespace, pod, container, this.http);
    pls.start();
    return pls;
  }
}
