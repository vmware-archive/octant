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
import { Streamer, ContentStreamService } from 'src/app/services/content-stream/content-stream.service';
import { IconService } from './services/icon.service';
import { ViewService } from './services/view/view.service';
import { BehaviorSubject } from 'rxjs';

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

  @ViewChild('scrollTarget') scrollTarget: ElementRef;

  hasTabs = false;
  hasReceivedContent = false;
  title: string = null;
  views: View[] = null;
  singleView: View = null;
  private iconName: string;

  constructor(
    private route: ActivatedRoute,
    private contentStreamService: ContentStreamService,
    private iconService: IconService,
    private viewService: ViewService
  ) {
    let streamer: Streamer = {
      behavior: this.behavior,
      handler: this.handleEvent,
    };
    this.contentStreamService.registerStreamer('content', streamer)
  }

  ngOnInit() {
    this.route.url.subscribe(url => {
      const currentPath = url.map(u => u.path).join('/');
      if (currentPath !== this.previousUrl) {
        this.title = null;
        this.singleView = null;
        this.views = null;
        this.previousUrl = currentPath;
        this.contentStreamService.openStream(currentPath);
        this.contentStreamService.streamer('content').subscribe(this.setContent);
        this.scrollTarget.nativeElement.scrollTop = 0;
      }
    });
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data) as ContentResponse;
    this.behavior.next(data);
    this.contentStreamService.removeAllSignals();
  };

  private setContent = (contentResponse: ContentResponse) => {
    const views = contentResponse.content.viewComponents;
    if (views.length === 0) {
      this.hasReceivedContent = false;
      return;
    }

    this.hasTabs = views.length > 1;
    if (this.hasTabs) {
      this.views = views;
      this.title = this.viewService.titleAsText(contentResponse.content.title);
    } else if (views.length === 1) {
      this.singleView = views[0];
    }

    this.hasReceivedContent = true;

    this.iconName = this.iconService.load(contentResponse.content);
  };

  ngOnDestroy() {
    this.contentStreamService.closeStream();
  }
}
