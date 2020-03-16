// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { NamespaceService } from 'src/app/modules/shared/services/namespace/namespace.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
})
export class NamespaceComponent implements OnInit, OnDestroy {
  namespaces: string[];
  currentNamespace = '';
  trackByIdentity = trackByIdentity;

  private namespaceSubscription: Subscription;

  constructor(private namespaceService: NamespaceService) {}

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
  }

  ngOnDestroy(): void {
    this.namespaceSubscription.unsubscribe();
  }

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }
}
