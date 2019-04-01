import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

import { ContentResponse } from '../../models/content';
import { Namespaces } from '../../models/namespace';
import { Navigation } from '../../models/navigation';
import { Filter, LabelFilterService } from '../label-filter/label-filter.service';

const pollEvery = 5;
const API_BASE = 'http://localhost:3001';

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

  constructor(private http: HttpClient, labelFilter: LabelFilterService) {
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
    this.eventSource = new EventSource(url);

    this.eventSource.addEventListener('message', (message: MessageEvent) => {
      const data = JSON.parse(message.data) as ContentResponse;
      this.content.next(data);
    });

    this.eventSource.addEventListener('navigation', (message: MessageEvent) => {
      const data = JSON.parse(message.data);
      this.navigation.next(data);
    });

    this.eventSource.addEventListener('namespaces', (message: MessageEvent) => {
      const data = JSON.parse(message.data) as Namespaces;
      this.namespaces.next(data.namespaces);
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

