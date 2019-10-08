// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import { PortForwardView } from 'src/app/models/content';

@Component({
  selector: 'app-view-port-forward',
  templateUrl: './port-forward.component.html',
  styleUrls: ['./port-forward.component.scss'],
})
export class PortForwardComponent {
  @Input() view: PortForwardView;
}
