// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import getAPIBase from '../../../../services/common/getAPIBase';
import {
  Streamer,
  ContentStreamService,
  ContextDescription,
} from '../../../../services/content-stream/content-stream.service';

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
  private behavior = new BehaviorSubject<KubeContextResponse>(emptyKubeContext);
  private contextsSource: BehaviorSubject<
    ContextDescription[]
  > = new BehaviorSubject<ContextDescription[]>([]);

  private selectedSource: BehaviorSubject<string> = new BehaviorSubject<string>(
    ''
  );

  constructor(
    private http: HttpClient,
    private contentStreamService: ContentStreamService
  ) {
    let streamer: Streamer = {
      behavior: this.behavior,
      handler: this.handleEvent,
    };
    this.contentStreamService.registerStreamer('kubeContext', streamer)

    contentStreamService.streamer('kubeContext').subscribe(update => {
      this.contextsSource.next(update.contexts);
      this.selectedSource.next(update.currentContext);
    });
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data) as KubeContextResponse;
    this.behavior.next(data);
  };

  select(context: ContextDescription) {
    this.selectedSource.next(context.name);

    this.updateContext(context.name).subscribe();
  }

  selected() {
    return this.selectedSource.asObservable();
  }

  contexts() {
    return this.contextsSource.asObservable();
  }

  private updateContext(name: string) {
    const url = [
      getAPIBase(),
      'api/v1/content/configuration',
      'kube-contexts',
    ].join('/');

    const payload = {
      requestedContext: name,
    };

    return this.http.post(url, payload);
  }
}
