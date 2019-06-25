// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, Input } from '@angular/core';

import { Navigation, NavigationChild } from '../../models/navigation';
import { IconService } from '../../modules/overview/services/icon.service';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent {
  @Input() navigation: Navigation = {
    sections: [],
  };

  constructor(private iconService: IconService) {}

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return this.iconService.load(item);
  }
}
