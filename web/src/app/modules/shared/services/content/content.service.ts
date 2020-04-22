/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Content, ContentResponse } from '../../models/content';
import { Params, Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../label-filter/label-filter.service';
import { NamespaceService } from '../namespace/namespace.service';

export const ContentUpdateMessage = 'event.octant.dev/content';

export interface ContentUpdate {
  content: Content;
  namespace: string;
  contentPath: string;
  queryParams: { [key: string]: string[] };
}

const emptyContentResponse: ContentResponse = {
  content: { extensionComponent: null, viewComponents: [], title: [] },
};

@Injectable({
  providedIn: 'root',
})
export class ContentService {
  defaultPath = new BehaviorSubject<string>('');
  current = new BehaviorSubject<ContentResponse>(emptyContentResponse);

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
    private namespaceService: NamespaceService
  ) {
    websocketService.registerHandler(ContentUpdateMessage, data => {
      const response = data as ContentUpdate;

      const s = JSON.stringify(data);
      if (s === this.lastReceived) {
        return;
      }

      this.lastReceived = s;

      this.setContent(response.content);
      namespaceService.setNamespace(response.namespace);

      if (response.contentPath) {
        if (this.previousContentPath.length > 0) {
          if (response.contentPath !== this.previousContentPath) {
            const segments = response.contentPath.split('/');
            this.router
              .navigate(segments, {
                queryParams: response.queryParams,
              })
              .catch(reason =>
                console.error(`unable to navigate`, { segments, reason })
              );
          }
        }
      }

      this.previousContentPath = response.contentPath;
    });

    labelFilterService.filters.subscribe(filters => {
      this.filters = filters;
    });
  }

  setContentPath(contentPath: string, params: Params) {
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

  private setContent(content: Content) {
    const contentResponse: ContentResponse = {
      content,
    };
    this.current.next(contentResponse);
  }
}
