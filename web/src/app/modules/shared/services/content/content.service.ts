/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';
import { BehaviorSubject, Observable } from 'rxjs';
import {
  Content,
  ContentResponse,
  LinkView,
  NamespacedTitle,
} from '../../models/content';
import { Params, Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../label-filter/label-filter.service';
import { NamespaceService } from '../namespace/namespace.service';
import { LoadingService } from '../loading/loading.service';
import { debounceTime, delay, distinctUntilChanged } from 'rxjs/operators';
import { Title } from '@angular/platform-browser';

export const ContentUpdateMessage = 'event.octant.dev/content';

export interface ContentUpdate {
  content: Content;
  namespace: string;
  contentPath: string;
  queryParams: { [key: string]: string[] };
}

const emptyContentResponse: ContentResponse = {
  content: { extensionComponent: null, viewComponents: [], title: [] },
  currentPath: '',
};

const emptyTitle: NamespacedTitle = {
  namespace: '',
  title: '',
  path: '',
};

@Injectable({
  providedIn: 'root',
})
export class ContentService {
  current = new BehaviorSubject<ContentResponse>(emptyContentResponse);
  title = new BehaviorSubject<NamespacedTitle>(emptyTitle);
  viewScrollPos = new BehaviorSubject<number>(0);
  debouncedScrollPos = new BehaviorSubject<number>(0);

  private previousContentPath = '';

  private filters: Filter[] = [];
  get currentFilters(): Filter[] {
    return this.filters;
  }

  private lastReceived = '';

  constructor(
    private router: Router,
    private websocketService: WebsocketService,
    private labelFilterService: LabelFilterService,
    private namespaceService: NamespaceService,
    private loadingService: LoadingService,
    private titleService: Title
  ) {
    websocketService.registerHandler(ContentUpdateMessage, data => {
      const response = data as ContentUpdate;

      const s = JSON.stringify(data);
      if (s === this.lastReceived) {
        return;
      }

      this.lastReceived = s;

      this.setContent(response);
      this.setTitle(response);
      namespaceService.setNamespace(response.namespace);
      if (response.contentPath) {
        if (this.previousContentPath.length > 0) {
          if (response.contentPath !== this.previousContentPath) {
            const segments = response.contentPath.split('/');
            this.router
              .navigate(segments, {
                queryParams: response.queryParams,
              })
              .then(result => {
                if (result) {
                  this.delayedComplete(true);
                } else {
                  this.loadingService.requestComplete.next(true);
                }
              })
              .catch(reason => {
                this.loadingService.requestComplete.next(true);
                console.error(`unable to navigate`, { segments, reason });
              });
          }
        } else {
          this.loadingService.requestComplete.next(true);
        }
      }

      this.previousContentPath = response.contentPath;
    });

    labelFilterService.filters.subscribe(filters => {
      this.filters = filters;
    });

    this.viewScrollPos
      .pipe(debounceTime(100), distinctUntilChanged())
      .subscribe(pos => this.debouncedScrollPos.next(pos));
  }

  delayedComplete(value: boolean) {
    const delayed = new Observable(x => {
      x.next();
    })
      .pipe(delay(700))
      .subscribe(() => {
        this.loadingService.requestComplete.next(value);
        delayed.unsubscribe();
      });
  }

  setContentPath(contentPath: string, params: Params) {
    this.viewScrollPos.next(0);
    if (contentPath === this.previousContentPath) {
      return;
    }

    if (!contentPath) {
      contentPath = '';
    }

    const payload = { contentPath, params };
    this.websocketService.sendMessage(
      'action.octant.dev/setContentPath',
      payload
    );
  }

  private setContent(contentUpdate: ContentUpdate) {
    const contentResponse: ContentResponse = {
      content: contentUpdate.content,
      currentPath: contentUpdate.contentPath,
    };
    this.current.next(contentResponse);
  }

  private setTitle(response: ContentUpdate) {
    const title = response?.content?.title;

    if (!title || title.length === 0) {
      return;
    }
    const titleView = title[title.length - 1] as LinkView;
    const titleVal = titleView.config.value;

    const pageTitle =
      response.namespace.length === 0
        ? `Octant | ${titleVal}`
        : `Octant | ${titleVal} | ${response.namespace}`;
    this.titleService.setTitle(pageTitle);
    this.title.next({
      namespace: response.namespace,
      title: titleVal,
      path: response.contentPath,
    });
  }

  setScrollPos(pos: number) {
    this.viewScrollPos.next(pos);
  }
}
