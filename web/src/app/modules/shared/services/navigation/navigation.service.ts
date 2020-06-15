/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';
import { NavigationEnd, Router, RouterEvent } from '@angular/router';
import { filter } from 'rxjs/operators';
import { LoadingService } from '../loading/loading.service';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  current = new BehaviorSubject<Navigation>(emptyNavigation);
  public lastSelection: BehaviorSubject<number> = new BehaviorSubject<number>(
    -1
  );
  public expandedState: BehaviorSubject<any> = new BehaviorSubject<any>({});
  public collapsed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
    true
  );
  activeUrl = new BehaviorSubject<string>('');

  constructor(
    private loadingService: LoadingService,
    private websocketService: WebsocketService,
    private contentService: ContentService,
    private router: Router
  ) {
    websocketService.registerHandler('event.octant.dev/navigation', data => {
      const update = data as Navigation;
      this.current.next(update);

      contentService.defaultPath.next(update.defaultPath);
      this.updateLastSelection();
    });

    router.events
      .pipe(filter(e => e instanceof NavigationEnd))
      .subscribe((event: RouterEvent) => {
        this.loadingService.requestComplete.next(false);
        this.activeUrl.next(event.url);
        this.updateLastSelection();
      });
  }

  updateLastSelection() {
    const targetUrl = this.activeUrl.value;
    let suggestedIndex = this.indexFromUrl(targetUrl);

    if (suggestedIndex === -1) {
      suggestedIndex = this.indexFromUrl(
        targetUrl.substring(0, targetUrl.lastIndexOf('/'))
      );
    }

    if (suggestedIndex >= 0 && suggestedIndex !== this.lastSelection.value) {
      this.lastSelection.next(suggestedIndex);
    }
  }

  indexFromUrl(url: string) {
    const short = url.substring(1);
    const paths = url.split('/');

    if (paths[1] === 'workloads') {
      return 0;
    }

    for (let index = 0; index < this.current.value.sections.length; index++) {
      const section = this.current.value.sections[index];

      if (section.path === short) {
        return index;
      } else if (section.children) {
        const suggested = section.children.findIndex(
          child => child.path === short
        );
        if (suggested >= 0) {
          return index;
        }
      }
    }
    return -1;
  }
}
