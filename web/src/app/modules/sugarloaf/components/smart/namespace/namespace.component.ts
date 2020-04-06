// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { NamespaceService } from 'src/app/modules/shared/services/namespace/namespace.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { Subscription } from 'rxjs';
import { NavigationService } from '../../../../shared/services/navigation/navigation.service';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
})
export class NamespaceComponent implements OnInit, OnDestroy {
  namespaces: string[];
  currentNamespace = '';
  trackByIdentity = trackByIdentity;
  navigation = {
    sections: [],
    defaultPath: '',
  };
  lastSelection: number;
  activeUrl: string;

  private namespaceSubscription: Subscription;

  constructor(
    private namespaceService: NamespaceService,
    private navigationService: NavigationService
  ) {}

  ngOnInit() {
    this.namespaceSubscription = this.namespaceService.activeNamespace.subscribe(
      (namespace: string) => {
        this.currentNamespace = namespace;
      }
    );

    this.namespaceSubscription = this.namespaceService.availableNamespaces.subscribe(
      (namespaces: string[]) => {
        this.namespaces = namespaces;
      }
    );

    this.navigationService.current.subscribe(
      navigation => (this.navigation = navigation)
    );

    this.navigationService.lastSelection.subscribe(
      selection => (this.lastSelection = selection)
    );

    this.navigationService.activeUrl.subscribe(url => (this.activeUrl = url));
  }

  ngOnDestroy(): void {
    this.namespaceSubscription.unsubscribe();
  }

  namespaceClass(namespace: string) {
    const active = this.currentNamespace === namespace ? ['active'] : [];
    return ['context-button', ...active];
  }

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }

  showDropdown() {
    if (this.lastSelection && this.navigation.sections[this.lastSelection]) {
      const url = this.navigation.sections[this.lastSelection].path;
      return !this.isClusterSpecific(url);
    }
    return true;
  }

  private namespaceFromUrl(): string {
    const paths = this.activeUrl.split('/');
    const len = paths.length;
    if (len > 1 && paths[len - 2] === 'namespaces') {
      const hashIndex = paths[len - 1].indexOf('#');
      return hashIndex > 0
        ? paths[len - 1].substring(0, hashIndex)
        : paths[len - 1];
    }
    return '';
  }

  private isClusterSpecific(url: string) {
    const ns = this.namespaceFromUrl();
    if (ns.length > 0) {
      this.selectNamespace(ns);
      return false;
    }

    if (url.includes('cluster-overview')) {
      return !this.activeUrl.endsWith(this.currentNamespace);
    }

    return false;
  }
}
