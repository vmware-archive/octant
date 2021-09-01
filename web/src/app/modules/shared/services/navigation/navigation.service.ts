/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Injectable } from '@angular/core';
import { WebsocketService } from '../../../../data/services/websocket/websocket.service';
import { BehaviorSubject } from 'rxjs';
import isSvg from 'is-svg';
import {
  Navigation,
  NavigationChild,
} from '../../../sugarloaf/models/navigation';
import { ContentService } from '../content/content.service';
import { NavigationEnd, Router, RouterEvent } from '@angular/router';
import { filter } from 'rxjs/operators';
import { LoadingService } from '../loading/loading.service';
import { ClarityIcons } from '@cds/core/icon';

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
  children?: NavigationChild[];
  customSvg?: string;
};

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

const MAX_RANK = 99;

@Injectable({
  providedIn: 'root',
})
export class NavigationService {
  current = new BehaviorSubject<Navigation>(emptyNavigation);
  modules = new BehaviorSubject<Module[]>([]);
  selectedItem = new BehaviorSubject<Selection>({ module: 0, index: -1 });
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
        const newUrl = update.defaultPath.startsWith('/')
          ? update.defaultPath
          : '/' + update.defaultPath;
        if (newUrl !== this.activeUrl.value) {
          this.activeUrl.next(newUrl);
        }
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
    const suggested = this.indexFromUrl(targetUrl);
    this.selectedItem.next(suggested);
  }

  indexFromUrl(url: string): Selection {
    const strippedUrl = this.stripUrl(url);
    let highestSelection = { module: 0, index: 0 };
    let rank = -1;
    let highestRank = -1;

    if (strippedUrl.length === 0) {
      return { module: 1, index: 0 };
    }

    for (const [moduleIndex, module] of this.modules.value.entries()) {
      const modulePath = this.stripUrl(module.path);

      rank = this.compareSubPaths(strippedUrl, modulePath);
      if (this.isExactMatch(rank)) {
        return { module: moduleIndex, index: 0 };
      } else if (rank > highestRank) {
        highestRank = rank;
        highestSelection = { module: moduleIndex, index: 0 };
      }

      for (const [childIndex, child] of module.children.entries()) {
        rank = this.compareSubPaths(strippedUrl, this.stripUrl(child.path));
        if (this.isExactMatch(rank)) {
          return { module: moduleIndex, index: childIndex };
        } else if (rank > highestRank) {
          highestRank = rank;
          highestSelection = { module: moduleIndex, index: childIndex };
        }
        if (child.children) {
          for (const grandchild of child.children) {
            rank = this.compareSubPaths(
              strippedUrl,
              this.stripUrl(grandchild.path)
            );
            if (this.isExactMatch(rank)) {
              return { module: moduleIndex, index: childIndex };
            } else if (rank > highestRank) {
              highestRank = rank;
              highestSelection = { module: moduleIndex, index: childIndex };
            }
          }
        }
      }
    }
    return highestSelection;
  }

  compareSubPaths(url: string, path: string): number {
    if (url === path) {
      // exact match
      return MAX_RANK;
    }

    let match = -1;
    const paths = path.split('/');
    const urls = url.split('/');

    do {
      match++;
    } while (paths[match] === urls[match]);

    return match;
  }

  isExactMatch(rank: number): boolean {
    return rank === MAX_RANK;
  }

  stripUrl(url: string): string {
    return url.startsWith('/') ? url.substring(1) : url;
  }

  createModules(sections: any[]) {
    const modules: Module[] = [];
    let pluginsIndex = 3;

    sections.forEach((section, index) => {
      if (section.module && section.module.length > 0) {
        if (section.customSvg) {
          this.registerCustomSvg(section);
        }

        modules.push({
          startIndex: index,
          name: section.module,
          icon: section.iconName,
          description: section.description,
          path: section.path,
          title: section.title,
          customSvg: section.customSvg,
        });
      }
    });

    modules.forEach((module, index) => {
      module.children = [];
      module.endIndex =
        index === modules.length - 1
          ? sections.length - 1
          : modules[index + 1].startIndex;
      if (module.name === 'configuration') {
        pluginsIndex = index;
      }
      if (sections[module.startIndex].children) {
        if (module.path !== sections[module.startIndex].children[0].path) {
          const first = {
            name: module.name,
            path: module.path,
            icon: module.icon,
            title: module.title,
            customSvg: module.customSvg,
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
      if (module.children) {
        this.findCustomSvg(module.children, true);
      }
    });
    if (modules.length > 0) {
      modules.push(modules.splice(pluginsIndex, 1)[0]);
    }
    this.modules.next(modules);
  }

  findCustomSvg(navigationChildren: NavigationChild[], skip: boolean): void {
    navigationChildren.forEach(navChild => {
      if (navChild.children) {
        this.findCustomSvg(navChild.children, false);
      }
      if (!skip && navChild.customSvg) {
        // do not show custom icon for parent of submenu
        this.registerCustomSvg(navChild);
      }
    });
  }

  registerCustomSvg(nc: NavigationChild): void {
    if (isSvg(nc.customSvg)) {
      ClarityIcons.addIcons([nc.iconName, nc.customSvg]);
    } else {
      console.error(
        `Invalid SVG for module: '${nc.title}'. Using default icon shape...`
      );
      nc.iconName = 'times';
    }
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
