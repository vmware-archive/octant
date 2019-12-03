import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Content, ContentResponse } from '../../../../models/content';
import { ActivatedRoute, Params, Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../../../../services/label-filter/label-filter.service';
import { NamespaceService } from '../../../../services/namespace/namespace.service';

export const ContentUpdateMessage = 'content';

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

  constructor(
    private router: Router,
    private websocketService: WebsocketService,
    private labelFilterService: LabelFilterService,
    private namespaceService: NamespaceService
  ) {
    websocketService.registerHandler(ContentUpdateMessage, data => {
      const response = data as ContentUpdate;
      this.setContent(response.content);
      namespaceService.setNamespace(response.namespace);

      if (response.contentPath) {
        if (this.previousContentPath.length > 0) {
          if (response.contentPath !== this.previousContentPath) {
            const segments = response.contentPath.split('/');
            this.router.navigate(segments, {
              queryParams: response.queryParams,
            });
          }
        }

        this.previousContentPath = response.contentPath;
      }
    });

    labelFilterService.filters.subscribe(filters => {
      this.filters = filters;
    });
  }

  setContentPath(contentPath: string, params: Params) {
    if (!contentPath) {
      contentPath = '';
    }

    const payload = { contentPath, params };
    this.websocketService.sendMessage('setContentPath', payload);
  }

  private setContent(content: Content) {
    const contentResponse: ContentResponse = {
      content,
    };
    this.current.next(contentResponse);
  }
}
