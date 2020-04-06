/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';
import { NavigationEnd, Router, RouterEvent } from '@angular/router';
import { filter } from 'rxjs/operators';

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
    private websocketService: WebsocketService,
    private contentService: ContentService,
    private router: Router
  ) {
    websocketService.registerHandler('event.octant.dev/navigation', data => {
      const update = data as Navigation;
      this.current.next(update);

      contentService.defaultPath.next(update.defaultPath);
    });

    router.events
      .pipe(filter(e => e instanceof NavigationEnd))
      .subscribe((event: RouterEvent) => {
        this.activeUrl.next(event.url);
      });
  }
}
