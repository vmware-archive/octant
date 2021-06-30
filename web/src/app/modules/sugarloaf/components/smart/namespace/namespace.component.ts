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
import {
  Module,
  NavigationService,
  Selection,
} from '../../../../shared/services/navigation/navigation.service';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NamespaceComponent implements OnInit, OnDestroy {
  readonly defaultNsLimit = 10;
  namespaces: string[];
  currentNamespace = '';
  trackByIdentity = trackByIdentity;
  modules: Module[] = [];
  selectedItem: Selection;
  activeUrl: string;
  showDropdown: boolean;
  nsLimit = this.defaultNsLimit;

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

    this.navigationService.modules.subscribe(modules => {
      this.modules = modules;
      this.cdr.detectChanges();
    });

    this.navigationService.selectedItem.subscribe(selection => {
      this.selectedItem = selection;
      this.showDropdown = this.hasDropdown();
      this.cdr.detectChanges();
    });
    this.navigationService.activeUrl.subscribe(url => {
      this.activeUrl = url;
      this.cdr.detectChanges();
    });
  }

  ngOnDestroy(): void {
    if (this.namespaceSubscription) {
      this.namespaceSubscription.unsubscribe();
    }
  }

  namespaceClass(namespace: string) {
    const active = this.currentNamespace === namespace ? ['active'] : [];
    return ['context-button', ...active];
  }

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }

  hasDropdown() {
    if (this.selectedItem && this.modules[this.selectedItem.module]) {
      return this.modules[this.selectedItem.module].name !== 'cluster-overview';
    }
    return true;
  }

  toggleShowMore(): void {
    this.nsLimit =
      this.nsLimit === this.namespaces.length
        ? this.defaultNsLimit
        : this.namespaces.length;
  }

  private routerLinkPath(namespace: string): string {
    return this.navigationService.redirect(namespace);
  }
}
