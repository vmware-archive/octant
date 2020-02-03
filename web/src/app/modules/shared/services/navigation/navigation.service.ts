/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  current = new BehaviorSubject<Navigation>(emptyNavigation);

  constructor(
    private websocketService: WebsocketService,
    private contentService: ContentService
  ) {
    websocketService.registerHandler('navigation', data => {
      const update = data as Navigation;
      this.current.next(update);

      contentService.defaultPath.next(update.defaultPath);
    });
  }
}
