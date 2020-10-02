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
import { isEqual } from 'lodash';
import { Subscription } from 'rxjs';
import { LoadingService } from 'src/app/modules/shared/services/loading/loading.service';

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
  private defaultPath: string;
  private previousParams: Params;
  private loadingSubscription: Subscription;
  public showSpinner = false;
  currentPath = '';

  constructor(
    private router: Router,
    private iconService: IconService,
    private contentService: ContentService,
    private loadingService: LoadingService
  ) {}

  ngOnInit() {
    this.updatePath(this.router.routerState.snapshot.url);

    this.contentSubscription = this.contentService.current.subscribe(
      contentResponse => {
        this.setContent(contentResponse);
      }
    );

    this.loadingService
      .withDelay(this.loadingService.requestComplete, 650, 1000)
      .subscribe(v => {
        this.showSpinner = v;
      });
  }

  ngOnDestroy() {
    this.resetView();
    if (this.contentSubscription) {
      this.contentSubscription.unsubscribe();
    }
    if (this.loadingSubscription) {
      this.loadingSubscription.unsubscribe();
    }
  }

  updatePath(url: string) {
    const tree = this.router.parseUrl(url);

    const primary = tree.root.children.primary;
    let segments = [];
    if (primary) {
      segments = primary.segments;
    }

    this.handlePathChange(segments, tree.queryParams, false);
  }

  private handlePathChange(
    segments: UrlSegment[],
    queryParams: Params,
    force: boolean
  ) {
    const urlPath = segments.map(u => u.path).join('/');
    this.currentPath = urlPath;
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
    this.views = null;
  }

  private setContent = (contentResponse: ContentResponse) => {
    if (
      this.currentPath.length > 0 &&
      contentResponse.currentPath !== this.currentPath
    ) {
      return; // ignore premature updates
    }

    const views = contentResponse.content.viewComponents;
    if (!views || views.length === 0) {
      this.hasReceivedContent = false;
      return;
    }

    this.buttonGroup = contentResponse.content.buttonGroup;

    this.extView = contentResponse.content.extensionComponent;
    this.views = views;
    this.title = contentResponse.content.title?.map((item: LinkView) => ({
      title: item.config.value,
      url: item.config.ref,
    }));

    this.hasReceivedContent = true;
  };

  onScroll(event) {
    this.contentService.setScrollPos(event.target.scrollTop);
  }
}
