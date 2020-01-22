// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { ActivatedRoute, Params, Router, UrlSegment } from '@angular/router';
import { ContentResponse, View, ExtensionView } from 'src/app/models/content';
import { IconService } from './services/icon.service';
import { ViewService } from './services/view/view.service';
import { BehaviorSubject, combineLatest } from 'rxjs';
import { untilDestroyed } from 'ngx-take-until-destroy';
import { ContentService } from './services/content/content.service';
import { WebsocketService } from './services/websocket/websocket.service';
import { KubeContextService } from './services/kube-context/kube-context.service';
import { take } from 'rxjs/operators';
import _ from 'lodash';

const emptyContentResponse: ContentResponse = {
  content: {
    extensionComponent: null,
    viewComponents: [],
    title: [],
  },
};

interface LocationCallbackOptions {
  segments: UrlSegment[];
  params: Params;
  currentContext: string;
  fragment: string;
}

@Component({
  selector: 'app-overview',
  templateUrl: './overview.component.html',
  styleUrls: ['./overview.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default,
})
export class OverviewComponent implements OnInit, OnDestroy {
  behavior = new BehaviorSubject<ContentResponse>(emptyContentResponse);
  @ViewChild('scrollTarget', { static: true }) scrollTarget: ElementRef;
  hasTabs = false;
  hasReceivedContent = false;
  title: string = null;
  views: View[] = null;
  extView: ExtensionView = null;
  singleView: View = null;
  private previousUrl = '';
  private iconName: string;
  private defaultPath: string;
  private previousParams: Params;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private iconService: IconService,
    private viewService: ViewService,
    private contentService: ContentService,
    private websocketService: WebsocketService,
    private kubeContextService: KubeContextService
  ) {}

  ngOnInit() {
    this.contentService.current
      .pipe(untilDestroyed(this))
      .subscribe(contentResponse => {
        this.setContent(contentResponse);
      });

    this.withCurrentLocation(options => {
      this.handlePathChange(options.segments, options.params, false);
    });

    this.websocketService.reconnected.subscribe(() => {
      this.withCurrentLocation(options => {
        this.handlePathChange(options.segments, options.params, true);
        this.kubeContextService.select({ name: options.currentContext });
      }, true);
    });
  }

  ngOnDestroy() {
    this.resetView();
  }

  private withCurrentLocation(
    callback: (options: LocationCallbackOptions) => void,
    takeOne = false
  ) {
    let observable = combineLatest([
      this.route.url,
      this.route.queryParams,
      this.route.fragment,
      this.kubeContextService.selected(),
    ]);

    if (takeOne) {
      observable = observable.pipe(take(1));
    }

    observable.subscribe(([segments, params, fragment, currentContext]) => {
      if (currentContext !== '') {
        callback({
          segments,
          params,
          fragment,
          currentContext,
        });
      }
    });
  }

  private handlePathChange(
    segments: UrlSegment[],
    queryParams: Params,
    force: boolean
  ) {
    const urlPath = segments.map(u => u.path).join('/');
    const currentPath = urlPath || this.defaultPath;
    if (
      force ||
      currentPath !== this.previousUrl ||
      !_.isEqual(queryParams, this.previousParams)
    ) {
      this.resetView();
      this.previousUrl = currentPath;
      this.previousParams = queryParams;
      this.contentService.setContentPath(currentPath, queryParams);
      this.scrollTarget.nativeElement.scrollTop = 0;
    }
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
      // TODO: show a loading screen here (#506)
      return;
    }

    const view = contentResponse.content.extensionComponent;
    this.extView = view;

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
}
