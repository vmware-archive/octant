// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Router } from '@angular/router';
import { WebsocketService } from '../../../modules/overview/services/websocket/websocket.service';

export interface Filter {
  key: string;
  value: string;
}

interface UpdateFilters {
  filters: Filter[];
}

@Injectable({
  providedIn: 'root',
})
export class LabelFilterService {
  public filters = new BehaviorSubject<Filter[]>([]);

  constructor(
    private router: Router,
    private websocketService: WebsocketService
  ) {
    websocketService.registerHandler('filters', data => {
      const update = data as UpdateFilters;
      this.filters.next(update.filters);
    });
  }

  add(filter: Filter): void {
    this.websocketService.sendMessage('addFilter', {
      filter,
    });
  }

  remove(filter: Filter): void {
    this.websocketService.sendMessage('removeFilter', {
      filter,
    });
  }

  clear(): void {
    this.websocketService.sendMessage('clearFilters', {});
  }

  decodeFilter(filterSource: string): Filter | null {
    const spl = filterSource.split(':');
    if (spl.length === 2) {
      return { key: spl[0], value: spl[1] };
    }
    return null;
  }
}
