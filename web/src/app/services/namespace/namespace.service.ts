// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { NavigationEnd, PRIMARY_OUTLET, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import _ from 'lodash';
import {
  Streamer,
  ContentStreamService,
} from '../content-stream/content-stream.service';
import { NavigationChild } from '../../models/navigation';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../notifier/notifier.service';
import { includesArray } from '../../util/includesArray';
import { Namespaces } from '../../models/namespace';
import getAPIBase from '../common/getAPIBase';

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  private notifierSession: NotifierSession;
  current = new BehaviorSubject<string>('default');
  list = new BehaviorSubject<string[]>([]);
  behavior = new BehaviorSubject<string[]>([]);

  constructor(
    private router: Router,
    private http: HttpClient,
    private contentStreamService: ContentStreamService,
    notifierService: NotifierService
  ) {
    this.notifierSession = notifierService.createSession();

    let streamer: Streamer = {
      behavior: this.behavior,
      handler: this.handleEvent,
    };
    this.contentStreamService.registerStreamer('namespaces', streamer);

    this.behavior.subscribe((namespaces: string[]) => {
      this.list.next(namespaces);
    });

    this.router.events.subscribe(event => {
      if (!(event instanceof NavigationEnd)) {
        return;
      }
      this.handleUrlPathChange();
    });
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data) as Namespaces;
    this.behavior.next(data.namespaces);
  };

  private handleUrlPathChange() {
    this.notifierSession.removeAllSignals();
    const namespace = this.getNamespaceFromUrlPath(this.router.url);

    if (namespace && !this.isNamespaceValid(namespace)) {
      this.notifierSession.pushSignal(
        NotifierSignalType.ERROR,
        'The current set namespace is not valid'
      );
      return;
    }

    const currentNS = this.current.getValue();
    if (currentNS !== namespace) {
      this.current.next(namespace);
    }
  }

  public getInitialNamespace() {
    const url = [getAPIBase(), 'api/v1/namespace'].join('/');
    return this.http.get(url);
  }

  private isNamespaceValid(namespaceToCheck: string): boolean {
    const listOfNamespaces = this.list.getValue();
    if (listOfNamespaces.length < 1) {
      // TODO: need a better way to check if
      // namespaces have loaded yet
      return true;
    }
    if (!namespaceToCheck) {
      return false;
    }
    return _.includes(listOfNamespaces, namespaceToCheck);
  }

  private getNamespaceFromUrlPath(url: string): string {
    if (!url) {
      throw new Error('No url');
    }
    const urlTree = this.router.parseUrl(url);
    const urlSegments = urlTree.root.children[PRIMARY_OUTLET].segments;
    if (urlSegments.length > 3 && urlSegments[2].path === 'namespace') {
      return urlSegments[3].path;
    }
  }

  private getPathArray(url: string): string[] {
    if (!url || url === '/') {
      return [];
    }
    const urlTree = this.router.parseUrl(url);
    return _.filter(
      _.map(urlTree.root.children[PRIMARY_OUTLET].segments, 'path')
    );
  }

  // When a user decides to switch namespaces, they expect to be able to
  // switch to the new namespace while still looking at the same resource.
  // Because our navigation is dynamic, this function traverses our navigation
  // options to match against the correct one.
  // See: issue #73
  private getNewRoute(namespace: string): string[] {
    const basePath = ['/content', 'overview', 'namespace'];
    const currentURLPath = this.getPathArray(this.router.url);

    // If the user is not on a namespace-scoped page but they switch
    // namespaces, take them to the namespace's overview page
    if (currentURLPath[2] !== 'namespace') {
      return [...basePath, namespace];
    }

    let routeCandidate = basePath;
    const navigation = this.contentStreamService
      .streamer('navigation')
      .getValue();

    navigationSectionLoop: for (const navigationSection of navigation.sections as NavigationChild[]) {
      const sectionPath = this.getPathArray(navigationSection.path);

      if (!includesArray(currentURLPath, sectionPath)) {
        continue;
      }

      let pointer = navigationSection;
      while (pointer && pointer.children) {
        const { children } = pointer;
        pointer = null;

        navigationChildLoop: for (const child of children as NavigationChild[]) {
          const pathArr = this.getPathArray(child.path);

          if (includesArray(currentURLPath, pathArr)) {
            routeCandidate = pathArr;
            // We're not guaranteed that namespace-scoped CRDs will exist
            // outside the current namespace so we move the user to the
            // CRD list view
            if (child.title !== 'Custom Resources') {
              pointer = child;
            }
            break navigationChildLoop;
          }
        }
      }

      break navigationSectionLoop;
    }

    // Set to new namespace
    routeCandidate[3] = namespace;
    return routeCandidate;
  }

  setNamespace(namespace: string) {
    this.current.next(namespace);
    const newRoute = this.getNewRoute(namespace);
    this.router.navigate(newRoute);
  }
}
