// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { NavigationEnd, PRIMARY_OUTLET, Router } from '@angular/router';
import { BehaviorSubject } from 'rxjs';
import _ from 'lodash';
import { ContentStreamService } from '../content-stream/content-stream.service';
import {
  NotifierService,
  NotifierSession,
  NotifierSignalType,
} from '../notifier/notifier.service';

@Injectable({
  providedIn: 'root',
})
export class NamespaceService {
  private notifierSession: NotifierSession;
  current = new BehaviorSubject<string>('default');
  list = new BehaviorSubject<string[]>([]);

  constructor(
    private router: Router,
    private contentStreamService: ContentStreamService,
    notifierService: NotifierService
  ) {
    this.notifierSession = notifierService.createSession();

    this.contentStreamService.namespaces.subscribe((namespaces: string[]) => {
      this.list.next(namespaces);
    });

    this.router.events.subscribe(event => {
      if (!(event instanceof NavigationEnd)) {
        return;
      }
      this.handleUrlPathChange();
    });
  }

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

  setNamespace(namespace: string) {
    this.current.next(namespace);
    this.router.navigate(['/content', 'overview', 'namespace', namespace]);
  }
}
