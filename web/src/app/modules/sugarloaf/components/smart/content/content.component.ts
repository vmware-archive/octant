/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */
import {
  ChangeDetectionStrategy,
  Component,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { Params, Router, UrlSegment } from '@angular/router';
import {
  ButtonGroupView,
  ContentResponse,
  ExtensionView,
  LinkView,
  PathItem,
  View,
} from 'src/app/modules/shared/models/content';
import { IconService } from '../../../../shared/services/icon/icon.service';
import { ContentService } from '../../../../shared/services/content/content.service';
import isEqual from 'lodash/isEqual';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-overview',
  templateUrl: './content.component.html',
  styleUrls: ['./content.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default,
})
export class ContentComponent implements OnInit, OnDestroy {
  hasTabs = false;
  hasReceivedContent = false;
  title: PathItem[] = null;
  views: View[] = null;
  extView: ExtensionView = null;
  singleView: View = null;
  buttonGroup: ButtonGroupView = null;
  private contentSubscription: Subscription;
  private previousUrl = '';
  private iconName: string;
  private defaultPath: string;
  private previousParams: Params;

  constructor(
    private router: Router,
    private iconService: IconService,
    private contentService: ContentService
  ) {}

  updatePath(url: string) {
    const tree = this.router.parseUrl(url);

    const primary = tree.root.children.primary;
    let segments = [];
    if (primary) {
      segments = primary.segments;
    }

    this.handlePathChange(segments, tree.queryParams, false);
  }

  ngOnInit() {
    this.updatePath(this.router.routerState.snapshot.url);

    this.contentSubscription = this.contentService.current.subscribe(
      contentResponse => {
        this.setContent(contentResponse);
      }
    );
  }

  ngOnDestroy() {
    this.resetView();
    this.contentSubscription.unsubscribe();
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
      (currentPath && currentPath !== this.previousUrl) ||
      !isEqual(queryParams, this.previousParams)
    ) {
      if (this.previousUrl === currentPath) {
        return;
      }

      this.previousParams = queryParams;
      this.resetView();
      this.contentService.setContentPath(currentPath, queryParams);
    }
  }

  private resetView() {
    this.title = null;
    this.singleView = null;
    this.views = null;
  }

  private setContent = (contentResponse: ContentResponse) => {
    const views = contentResponse.content.viewComponents;
    if (!views || views.length === 0) {
      this.hasReceivedContent = false;
      // TODO: show a loading screen here (#506)
      return;
    }
    this.buttonGroup = contentResponse.content.buttonGroup;

    this.extView = contentResponse.content.extensionComponent;
    this.hasTabs = views.length > 1;
    if (this.hasTabs) {
      this.views = views;
      this.title = contentResponse.content.title.map((item: LinkView) => ({
        title: item.config.value,
        url: item.config.ref,
      }));
    } else if (views.length === 1) {
      this.views = null;
      this.singleView = views[0];
    }

    this.hasReceivedContent = true;
    this.iconName = this.iconService.load(contentResponse.content);
  };
}
