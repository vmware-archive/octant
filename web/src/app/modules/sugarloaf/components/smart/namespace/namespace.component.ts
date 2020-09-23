// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { NamespaceService } from 'src/app/modules/shared/services/namespace/namespace.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { Subscription } from 'rxjs';
import { NavigationService } from '../../../../shared/services/navigation/navigation.service';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
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
  showDropdown: boolean;

  private namespaceSubscription: Subscription;

  constructor(
    private namespaceService: NamespaceService,
    private navigationService: NavigationService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.namespaceSubscription = this.namespaceService.activeNamespace.subscribe(
      (namespace: string) => {
        this.currentNamespace = namespace;
        this.cdr.detectChanges();
      }
    );

    this.namespaceSubscription = this.namespaceService.availableNamespaces.subscribe(
      (namespaces: string[]) => {
        this.namespaces = namespaces;
        this.cdr.detectChanges();
      }
    );

    this.navigationService.current.subscribe(navigation => {
      this.navigation = navigation;
      this.cdr.detectChanges();
    });

    this.navigationService.activeUrl.subscribe(url => {
      this.activeUrl = url;
      this.cdr.detectChanges();
    });

    this.navigationService.lastSelection.subscribe(selection => {
      this.lastSelection = selection;
      this.showDropdown = this.hasDropdown();
      this.cdr.detectChanges();
    });
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

  hasDropdown(): boolean {
    if (this.lastSelection && this.navigation.sections[this.lastSelection]) {
      const url = this.navigation.sections[this.lastSelection].path;
      return !this.isClusterSpecific(url);
    }
    return true;
  }

  private namespaceFromUrl(url: string): string {
    if (url) {
      const paths = url.split('/');
      const len = paths.length;
      if (len > 1 && paths[len - 2] === 'namespaces') {
        const hashIndex = paths[len - 1].indexOf('#');
        return hashIndex > 0
          ? paths[len - 1].substring(0, hashIndex)
          : paths[len - 1];
      }
    }
    return '';
  }

  private isClusterSpecific(url: string) {
    const ns = this.namespaceFromUrl(url);
    if (ns.length > 0) {
      this.selectNamespace(ns);
      return false;
    }

    if (url.includes('cluster-overview')) {
      return (
        this.currentNamespace.length === 0 ||
        !url.endsWith(this.currentNamespace)
      );
    }

    return false;
  }
}
