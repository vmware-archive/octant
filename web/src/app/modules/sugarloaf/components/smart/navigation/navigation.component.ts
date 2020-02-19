// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component, OnDestroy, OnInit } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Navigation, NavigationChild } from '../../../models/navigation';
import { IconService } from '../../../../shared/services/icon/icon.service';
import { NavigationService } from '../../../../shared/services/navigation/navigation.service';
import { untilDestroyed } from 'ngx-take-until-destroy';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent implements OnInit, OnDestroy {
  behavior = new BehaviorSubject<Navigation>(emptyNavigation);
  collapsed = false;
  navExpandedState: any;

  navigation = emptyNavigation;

  constructor(
    private iconService: IconService,
    private navigationService: NavigationService
  ) {}

  ngOnInit() {
    this.navigationService.current
      .pipe(untilDestroyed(this))
      .subscribe(navigation => (this.navigation = navigation));
    this.navigationService.expandedState
      .pipe(untilDestroyed(this))
      .subscribe(expanded => (this.navExpandedState = expanded));
  }

  ngOnDestroy() {}

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return this.iconService.load(item);
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data);
    this.behavior.next(data);
  };

  formatPath(path: string): string {
    if (!path.startsWith('/')) {
      return '/' + path;
    }

    return path;
  }

  setNavState($event, state: number) {
    this.navExpandedState[state] = $event;
    this.navigationService.expandedState.next(this.navExpandedState);
  }

  shouldExpand(index: number) {
    if (index.toString() in this.navExpandedState) {
      return this.navExpandedState[index];
    }
    return false;
  }
}
