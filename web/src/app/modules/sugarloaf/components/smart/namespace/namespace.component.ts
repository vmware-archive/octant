// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { untilDestroyed } from 'ngx-take-until-destroy';
import { NamespaceService } from 'src/app/modules/shared/services/namespace/namespace.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
})
export class NamespaceComponent implements OnInit, OnDestroy {
  namespaces: string[];
  currentNamespace = '';
  trackByIdentity = trackByIdentity;

  constructor(private namespaceService: NamespaceService) {}

  ngOnInit() {
    this.namespaceService.activeNamespace
      .pipe(untilDestroyed(this))
      .subscribe((namespace: string) => {
        this.currentNamespace = namespace;
      });

    this.namespaceService.availableNamespaces
      .pipe(untilDestroyed(this))
      .subscribe((namespaces: string[]) => {
        this.namespaces = namespaces;
      });
  }

  ngOnDestroy() {}

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }
}
