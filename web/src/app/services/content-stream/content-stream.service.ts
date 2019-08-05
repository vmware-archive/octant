// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { Location } from '@angular/common';
import { BehaviorSubject } from 'rxjs';
import getAPIBase from '../common/getAPIBase';
import { ContentResponse } from '../../models/content';

import {
  Filter,
  LabelFilterService,
} from '../label-filter/label-filter.service';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../notifier/notifier.service';
import { EventSourceService } from './event-source.service';
import _ from 'lodash';

export interface Streamer {
  handler: EventListenerOrEventListenerObject;
  behavior: BehaviorSubject<any>;
}

export interface ContextDescription {
  name: string;
}

const pollEvery = 5;
const API_BASE = getAPIBase();


@Injectable({
  providedIn: 'root',
})
export class ContentStreamService {

  private eventSource: EventSource;
  private notifierSession: NotifierSession;
  private currentPath: string;
  private streamers = new Map<string, Streamer>();

  constructor(
    private notifierService: NotifierService,
    private location: Location,
    private eventSourceService: EventSourceService,
    private labelFilterService: LabelFilterService
  ) {
    this.labelFilterService.filters.subscribe(() => this.restartStream());
    this.notifierSession = this.notifierService.createSession();
  }

  openStream(path: string) {
    this.closeStream();
    this.currentPath = path;
    const eventSourceUrl = this.createEventSourceUrl(path);
    this.eventSource = this.eventSourceService.createEventSource(
      eventSourceUrl
    );
    this.notifierSession.pushSignal(NotifierSignalType.LOADING, true);
    this.eventSource.addEventListener('error', this.handleErrorEvent);
    this.eventSource.addEventListener('objectNotFound', this.handleObjectNotFoundEvent);

    let eventSource = this.eventSource;
    this.streamers.forEach(function(value, key) {
      eventSource.addEventListener(key, value.handler);
    });
  }

  streamer(name: string): BehaviorSubject<any> {
    return this.streamers.get(name).behavior;
  }

  registerStreamer(name: string, streamer: Streamer) {
    this.streamers.set(name, streamer);
  }

  removeAllSignals() {
    this.notifierSession.removeAllSignals();
  }

  closeStream() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    this.currentPath = null;
    this.removeAllSignals();
  }

  private handleObjectNotFoundEvent = (message: MessageEvent) => {
    const redirectPath = message.data as string;
    this.location.go(redirectPath);
    this.currentPath = redirectPath.replace(/^(\/content\/)/, '');
    this.restartStream();
    this.notifierSession.pushSignal(
      NotifierSignalType.WARNING,
      'Kubernetes object was deleted from the cluster.'
    );
  };

  private handleErrorEvent = () => {
    this.notifierSession.pushSignal(
      NotifierSignalType.ERROR,
      'Lost back end source. Currently retrying...'
    );
  };

  private createEventSourceUrl = (path: string) => {
    const filters = this.labelFilterService.filters.getValue();

    let filterQuery = _.reduce(
      filters,
      (prev: string, cur: Filter, i: number) => {
        return (
          prev +
          (i > 0 ? '&' : '') +
          'filter=' +
          encodeURIComponent(`${cur.key}:${cur.value}`)
        );
      },
      ''
    );

    if (filterQuery.length > 0) {
      filterQuery = `&${filterQuery}`;
    }

    if (_.last(path) !== '/') {
      path += '/';
    }

    return `${API_BASE}/api/v1/content/${path}?poll=${pollEvery}${filterQuery}`;
  };

  private restartStream() {
    if (this.currentPath) {
      const path = this.currentPath;
      this.closeStream();
      this.openStream(path);
    }
  }
}
