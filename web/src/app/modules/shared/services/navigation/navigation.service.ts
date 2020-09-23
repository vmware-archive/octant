/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import { Navigation } from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';
import { NavigationEnd, Router, RouterEvent } from '@angular/router';
import { filter } from 'rxjs/operators';
import { LoadingService } from '../loading/loading.service';

export type Selection = {
  module: number;
  index: number;
};

export type Module = {
  name: string;
  title?: string;
  path?: string;
  description: string;
  startIndex: number;
  endIndex?: number;
  icon: string;
  children?: any[];
};

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  current = new BehaviorSubject<Navigation>(emptyNavigation);
  modules = new BehaviorSubject<Module[]>([]);
  selectedItem = new BehaviorSubject<Selection>({ module: 0, index: -1 });
  public expandedState: BehaviorSubject<any> = new BehaviorSubject<any>({});
  public collapsed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
    false
  );
  public showLabels: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(
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
      this.createModules(update.sections);
      if (update.defaultPath) {
        this.activeUrl.next(update.defaultPath);
      }
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
    let suggested = this.indexFromUrl(targetUrl);

    if (suggested.index === -1) {
      suggested = this.indexFromUrl(
        targetUrl.substring(0, targetUrl.lastIndexOf('/'))
      );
    }

    if (
      (suggested.index >= 0 &&
        suggested.module !== this.selectedItem.value.module) ||
      suggested.index !== this.selectedItem.value.index
    ) {
      this.selectedItem.next(suggested);
    }
  }

  indexFromUrl(url: string): Selection {
    const strippedUrl = this.stripUrl(url);
    for (const [moduleIndex, module] of this.modules.value.entries()) {
      const modulePath = this.stripUrl(module.path);

      if (strippedUrl === modulePath) {
        return { module: moduleIndex, index: 0 };
      } else {
        for (const [childIndex, child] of module.children.entries()) {
          if (strippedUrl === child.path) {
            return { module: moduleIndex, index: childIndex };
          }
          if (child.children) {
            for (const grandchild of child.children) {
              if (strippedUrl === grandchild.path) {
                this.expandedState[childIndex] = true;
                this.expandedState.next(this.expandedState);
                return { module: moduleIndex, index: childIndex };
              }
            }
          }
        }
      }
    }
    return { module: 0, index: -1 };
  }

  stripUrl(url: string): string {
    return url.startsWith('/') ? url.substring(1) : url;
  }

  createModules(sections: any[]) {
    const modules: Module[] = [];

    sections.forEach((section, index) => {
      if (section.module && section.module.length > 0) {
        modules.push({
          startIndex: index,
          name: section.module,
          icon: section.iconName,
          description: section.description,
          path: section.path,
          title: section.title,
        });
      }
    });
    modules.forEach((module, index) => {
      module.children = [];
      module.endIndex =
        index === modules.length - 1
          ? sections.length - 1
          : modules[index + 1].startIndex;
      if (sections[module.startIndex].children) {
        if (module.path !== sections[module.startIndex].children[0].path) {
          const first = {
            name: module.name,
            path: module.path,
            icon: module.icon,
            title: module.title,
          };
          module.children = [
            ...[first],
            ...sections[module.startIndex].children,
          ];
        } else {
          module.children = sections[module.startIndex].children;
        }
      } else {
        for (let i = module.startIndex; i < module.endIndex; i++) {
          module.children.push(sections[i]);
        }
      }
    });
    this.modules.next(modules);
  }

  redirect(namespace: string): string {
    let routerLink = '';
    const paths = this.activeUrl.value.split('/');
    const module = paths[1];

    switch (module) {
      case 'workloads': {
        routerLink = '/workloads/namespace/' + namespace;
        break;
      }
      case 'overview': {
        if (paths.length > 6) {
          routerLink = '/overview/namespace/' + namespace;
        } else {
          paths[3] = namespace;
          routerLink = paths.join('/');
        }
        break;
      }
      default: {
        routerLink = this.activeUrl.value;
      }
    }
    return routerLink;
  }
}
