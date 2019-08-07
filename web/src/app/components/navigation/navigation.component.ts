// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import {
  Streamer,
  ContentStreamService,
} from '../../services/content-stream/content-stream.service';
import { Navigation, NavigationChild } from '../../models/navigation';
import { IconService } from '../../modules/overview/services/icon.service';

const emptyNavigation: Navigation = {
  sections: [],
};

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent {
  behavior = new BehaviorSubject<Navigation>(emptyNavigation);

  @Input() navigation: Navigation = {
    sections: [],
  };

  constructor(
    private iconService: IconService,
    private contentStreamService: ContentStreamService
  ) {
    let streamer: Streamer = {
      behavior: this.behavior,
      handler: this.handleEvent,
    };

    this.contentStreamService.registerStreamer('navigation', streamer);
  }

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return this.iconService.load(item);
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data);
    this.behavior.next(data);
  };
}
