/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Content, ContentResponse } from '../../../../models/content';
import { Params, Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../../../../services/label-filter/label-filter.service';

export const ContentUpdateMessage = 'content';
export const ContentPathUpdateMessage = 'contentPath';

export interface ContentPathUpdate {
  contentPath: string;
  queryParams: { [key: string]: string[] };
}

const emptyContentResponse: ContentResponse = {
  content: { viewComponents: [], title: [] },
};

@Injectable({
  providedIn: 'root',
})
export class ContentService {
  defaultPath = new BehaviorSubject<string>('');
  current = new BehaviorSubject<ContentResponse>(emptyContentResponse);

  private filters: Filter[] = [];
  get currentFilters(): Filter[] {
    return this.filters;
  }

  constructor(
    private router: Router,
    private websocketService: WebsocketService,
    private labelFilterService: LabelFilterService
  ) {
    websocketService.registerHandler(ContentUpdateMessage, data => {
      const content = data as Content;
      this.setContent(content);
    });
    websocketService.registerHandler(ContentPathUpdateMessage, data => {
      const contentPathUpdate = data as ContentPathUpdate;
      this.router.navigate(
        ['content', ...contentPathUpdate.contentPath.split('/')],
        {
          queryParams: contentPathUpdate.queryParams,
        }
      );
    });

    labelFilterService.filters.subscribe(filters => {
      this.filters = filters;
    });
  }

  setContentPath(contentPath: string) {
    this.websocketService.sendMessage('setContentPath', {
      contentPath,
      filters: this.filters,
    });
  }

  setQueryParams(params: Params) {
    this.websocketService.sendMessage('setQueryParams', {
      params,
    });
  }

  private setContent(content: Content) {
    const contentResponse: ContentResponse = {
      content,
    };
    this.current.next(contentResponse);
  }
}
