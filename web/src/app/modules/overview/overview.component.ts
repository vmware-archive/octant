// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ContentResponse, View } from 'src/app/models/content';
import { IconService } from './services/icon.service';
import { ViewService } from './services/view/view.service';
import { BehaviorSubject } from 'rxjs';
import { ContentService } from './services/content/content.service';
import { WebsocketService } from './services/websocket/websocket.service';

const emptyContentResponse: ContentResponse = {
  content: {
    viewComponents: [],
    title: [],
  },
};

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.scss'],
})
export class OverviewComponent implements OnInit, OnDestroy {
  behavior = new BehaviorSubject<ContentResponse>(emptyContentResponse);

  private previousUrl = '';

  @ViewChild('scrollTarget', { static: true }) scrollTarget: ElementRef;

  hasTabs = false;
  hasReceivedContent = false;
  title: string = null;
  views: View[] = null;
  singleView: View = null;
  private iconName: string;
  private defaultPath: string;

  constructor(
    private route: ActivatedRoute,
    private iconService: IconService,
    private viewService: ViewService,
    private contentService: ContentService,
    private websocketService: WebsocketService
  ) {
    this.contentService.current.subscribe(contentResponse => {
      this.setContent(contentResponse);
    });
  }

  ngOnInit() {
    this.route.url.subscribe(this.handlePathChange());
    this.route.queryParams.subscribe(queryParams =>
      this.contentService.setQueryParams(queryParams)
    );

    this.websocketService.reconnected.subscribe(_ => {
      // when reconnecting, ensure the backend knows our path
      this.route.url.subscribe(this.handlePathChange(true)).unsubscribe();
      this.route.queryParams
        .subscribe(queryParams =>
          this.contentService.setQueryParams(queryParams)
        )
        .unsubscribe();
    });
  }

  private handlePathChange(force = false) {
    return url => {
      const currentPath = url.map(u => u.path).join('/') || this.defaultPath;
      if (currentPath !== this.previousUrl || force) {
        this.resetView();
        this.previousUrl = currentPath;
        this.contentService.setContentPath(currentPath);
        this.scrollTarget.nativeElement.scrollTop = 0;
      }
    };
  }

  private resetView() {
    this.title = null;
    this.singleView = null;
    this.views = null;
    this.hasReceivedContent = false;
  }

  private setContent = (contentResponse: ContentResponse) => {
    const views = contentResponse.content.viewComponents;
    if (!views || views.length === 0) {
      this.hasReceivedContent = false;
      // TODO: show a loading screen here
      return;
    }

    this.hasTabs = views.length > 1;
    if (this.hasTabs) {
      this.views = views;
      this.title = this.viewService.titleAsText(contentResponse.content.title);
    } else if (views.length === 1) {
      this.views = null;
      this.singleView = views[0];
    }

    this.hasReceivedContent = true;
    this.iconName = this.iconService.load(contentResponse.content);
  };

  ngOnDestroy() {
    this.resetView();
  }
}
