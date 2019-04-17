import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import {Location} from '@angular/common';
import { BehaviorSubject, Observable } from 'rxjs';
import getAPIBase from '../common/getAPIBase';
import { ContentResponse } from '../../models/content';
import { Namespaces } from '../../models/namespace';
import { Navigation } from '../../models/navigation';
import { Filter, LabelFilterService } from '../label-filter/label-filter.service';
import { NotifierService } from '../notifier/notifier.service';

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

  private content = new BehaviorSubject<ContentResponse>(emptyContentResponse);
  private namespaces = new BehaviorSubject<string[]>([]);
  private navigation = new BehaviorSubject<Navigation>(emptyNavigation);

  private filters: Filter[] = [];
  private currentPath: string;

  constructor(private http: HttpClient,
              labelFilter: LabelFilterService,
              private notifierService: NotifierService,
              private location: Location) {
    labelFilter.filters.subscribe((filters) => {
      this.filters = filters;
      this.restartPoller();
    });
  }

  getNavigation() {
    return this.http.get(`${API_BASE}/api/v1/navigation`);
  }

  getNamespaces() {
    return this.http.get(`${API_BASE}/api/v1/namespaces`);
  }

  private restartPoller() {
    if (this.currentPath) {
      const path = this.currentPath;

      this.stopPoller();
      this.startPoller(path);
    }
  }

  startPoller(path: string) {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }

    // if path ends with a namespace and no slash, append a slash
    if (path.match(/namespace\/.*[^\/]$/)) {
      path = path + '/';
    }

    this.currentPath = path;

    const filters = this.filters;

    let filterQuery = filters.reduce((prev: string, cur: Filter, i: number) => {
      return prev + (i > 0 ? '&' : '') + 'filter=' + encodeURIComponent(`${cur.key}:${cur.value}`);
    }, '');
    if (filterQuery.length > 0) {
      filterQuery = `&${filterQuery}`;
    }

    const url = `${API_BASE}/api/v1/content/${path}?poll=${pollEvery}${filterQuery}`;
    this.notifierService.loading.next(true);
    this.eventSource = new EventSource(url);

    this.eventSource.addEventListener('content', (message: MessageEvent) => {
      const data = JSON.parse(message.data) as ContentResponse;
      this.content.next(data);
      this.notifierService.loading.next(false);
      this.notifierService.error.next(null);
    });

    this.eventSource.addEventListener('navigation', (message: MessageEvent) => {
      const data = JSON.parse(message.data);
      this.navigation.next(data);
    });

    this.eventSource.addEventListener('namespaces', (message: MessageEvent) => {
      const data = JSON.parse(message.data) as Namespaces;
      this.namespaces.next(data.namespaces);
    });

    this.eventSource.addEventListener('objectNotFound', (message: MessageEvent) => {
      const redirectPath = message.data as string;
      this.location.go(redirectPath);
      this.currentPath = redirectPath.replace(/^(\/content\/)/, '');
      this.restartPoller();
      this.notifierService.warning.next('Kubernetes object was deleted from the cluster.');
    });

    this.eventSource.addEventListener('error', () => {
      this.notifierService.error.next('Lost back end source. Currently retrying...');
    });
  }

  pollNavigation(): Observable<Navigation> {
    return this.navigation;
  }

  pollContent(): Observable<ContentResponse> {
    return this.content;
  }

  pollNamespaces(): Observable<string[]> {
    return this.namespaces;
  }

  stopPoller() {
    if (this.eventSource) {
      this.eventSource.close();
    }

    this.currentPath = undefined;
  }
}
