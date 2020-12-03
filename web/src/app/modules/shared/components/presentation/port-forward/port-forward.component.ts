// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { PortForwardView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-port-forward',
  templateUrl: './port-forward.component.html',
  styleUrls: ['./port-forward.component.scss'],
})
export class PortForwardComponent extends AbstractViewComponent<PortForwardView> {
  constructor() {
    super();
  }

  update() {}
}
