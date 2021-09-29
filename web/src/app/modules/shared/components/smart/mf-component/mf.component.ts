/*
 * Copyright (c) 2021 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { Component } from '@angular/core';
import { MfComponentView } from '../../../models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { WebComponentWrapperOptions } from '@angular-architects/module-federation-tools';

@Component({
  selector: 'app-mf-component',
  templateUrl: './mf.component.html',
  styleUrls: ['./mf.component.scss'],
})
export class MfComponent extends AbstractViewComponent<MfComponentView> {
  name: string;
  meta: WebComponentWrapperOptions;

  constructor() {
    super();
  }

  update() {
    const config = this.v?.config;
    if (config) {
      this.name = config.name;
      this.meta = {
        remoteEntry: config.remoteEntry,
        remoteName: config.remoteName,
        exposedModule: config.exposedModule,
        elementName: config.elementName,
      };
    }
  }
}
