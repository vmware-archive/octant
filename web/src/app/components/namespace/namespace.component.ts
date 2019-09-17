// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnInit } from '@angular/core';
import { NamespaceService } from 'src/app/services/namespace/namespace.service';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.scss'],
})
export class NamespaceComponent implements OnInit {
  namespaces: string[];
  currentNamespace = '';
  trackByIdentity = trackByIdentity;

  constructor(private namespaceService: NamespaceService) {}

  ngOnInit() {
    this.namespaceService.activeNamespace.subscribe((namespace: string) => {
      this.currentNamespace = namespace;
    });

    this.namespaceService.availableNamespaces.subscribe(
      (namespaces: string[]) => {
        this.namespaces = namespaces;
      }
    );
  }

  selectNamespace(namespace: string) {
    this.namespaceService.setNamespace(namespace);
  }
}
