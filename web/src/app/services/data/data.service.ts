import { Injectable } from '@angular/core';
import {Location} from '@angular/common';
import { BehaviorSubject, Observable } from 'rxjs';
import getAPIBase from '../common/getAPIBase';
import { ContentResponse } from '../../models/content';
import { Namespaces } from '../../models/namespace';
import { Navigation } from '../../models/navigation';
import { Filter, LabelFilterService } from '../label-filter/label-filter.service';
import { NotifierService, NotifierSession, NotifierSignalType } from '../notifier/notifier.service';
import { EventSourceService } from './event-source.service';
import _ from 'lodash';

const pollEvery = 5;
const API_BASE = getAPIBase();

const emptyContentResponse: ContentResponse = {
  content: {
    viewComponents: [],
    title: [],
  },
};

const emptyNavigation: Navigation = {
  sections: [],
};

@Injectable({
  providedIn: 'root',
})
export class DataService {
  private eventSource: EventSource;
  private notifierSession: NotifierSession;
  private currentPath: string;

  content = new BehaviorSubject<ContentResponse>(emptyContentResponse);
  namespaces = new BehaviorSubject<string[]>([]);
  navigation = new BehaviorSubject<Navigation>(emptyNavigation);

  constructor(
    private notifierService: NotifierService,
    private location: Location,
    private eventSourceService: EventSourceService,
    private labelFilterService: LabelFilterService
  ) {
    this.labelFilterService.filters.subscribe(() => this.restartStream());
    this.notifierSession = this.notifierService.createSession();
  }

  private handleContentEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data) as ContentResponse;
    this.content.next(data);
    this.notifierSession.removeAllSignals();
  }

  private handleNavigationEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data);
    this.navigation.next(data);
  }

  private handleNamespaceEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data) as Namespaces;
    this.namespaces.next(data.namespaces);
  }

  private handleObjectNotFoundEvent = (message: MessageEvent) => {
    const redirectPath = message.data as string;
    this.location.go(redirectPath);
    this.currentPath = redirectPath.replace(/^(\/content\/)/, '');
    this.restartStream();
    this.notifierSession.pushSignal(NotifierSignalType.WARNING, 'Kubernetes object was deleted from the cluster.');
  }

  private handleErrorEvent = () => {
    this.notifierSession.pushSignal(NotifierSignalType.ERROR, 'Lost back end source. Currently retrying...');
  }

  private createEventSourceUrl = (path: string) => {
    const filters = this.labelFilterService.filters.getValue();

    let filterQuery = _.reduce(filters, (prev: string, cur: Filter, i: number) => {
      return prev + (i > 0 ? '&' : '') + 'filter=' + encodeURIComponent(`${cur.key}:${cur.value}`);
    }, '');

    if (filterQuery.length > 0) {
      filterQuery = `&${filterQuery}`;
    }

    if (_.last(path) !== '/') {
      path += '/';
    }

    return `${API_BASE}/api/v1/content/${path}?poll=${pollEvery}${filterQuery}`;
  }

  private restartStream() {
    if (this.currentPath) {
      const path = this.currentPath;
      this.closeStream();
      this.openStream(path);
    }
  }

  openStream(path: string) {
    this.closeStream();
    this.currentPath = path;
    const eventSourceUrl = this.createEventSourceUrl(path);
    this.eventSource = this.eventSourceService.createEventSource(eventSourceUrl);
    this.notifierSession.pushSignal(NotifierSignalType.LOADING, true);
    this.eventSource.addEventListener('content', this.handleContentEvent);
    this.eventSource.addEventListener('navigation', this.handleNavigationEvent);
    this.eventSource.addEventListener('namespaces', this.handleNamespaceEvent);
    this.eventSource.addEventListener('error', this.handleErrorEvent);
    this.eventSource.addEventListener('objectNotFound', this.handleObjectNotFoundEvent);
  }

  closeStream() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    this.currentPath = null;
    this.notifierSession.removeAllSignals();
  }
}
