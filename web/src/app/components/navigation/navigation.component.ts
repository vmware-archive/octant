// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnInit } from '@angular/core';
import { Navigation, NavigationChild } from '../../models/navigation';
import { IconService } from '../../modules/overview/services/icon.service';
import { NavigationService } from '../../modules/overview/services/navigation/navigation.service';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent implements OnInit {
  navigation = emptyNavigation;
  activeTab: string;

  constructor(
    private iconService: IconService,
    private navigationService: NavigationService,
    private router: Router
  ) {
    this.navigationService.current.subscribe(navigation => {
      this.navigation = navigation;
      this.updateNavigation(this.router.url);
    });

    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        const url = (event as NavigationEnd).url;
        this.updateNavigation(url);
      }
    });
  }

  ngOnInit(): void {}

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return this.iconService.load(item);
  }

  clickTab(tabName: string) {
    if (this.activeTab === tabName) {
      return;
    }

    this.activeTab = tabName;
  }

  private updateNavigation(url: string) {
    const path = url.replace(/^\//, '');
    this.navigation.sections.forEach(section => {
      if (path.startsWith(section.path)) {
        this.activeTab = section.title;
      }
    });
  }
}
