// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input } from '@angular/core';
import { View } from 'src/app/modules/shared/models/content';

@Component({
  selector: 'app-content-switcher',
  templateUrl: './content-switcher.component.html',
  styleUrls: ['./content-switcher.component.scss'],
})
export class ContentSwitcherComponent {
  @Input() view: View;
  constructor() {}
}
