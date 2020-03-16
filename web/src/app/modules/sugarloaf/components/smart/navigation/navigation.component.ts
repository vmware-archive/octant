// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { BehaviorSubject, Subscription } from 'rxjs';
import { Navigation, NavigationChild } from '../../../models/navigation';
import { IconService } from '../../../../shared/services/icon/icon.service';
import { NavigationService } from '../../../../shared/services/navigation/navigation.service';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavigationComponent implements OnInit, OnDestroy {
  behavior = new BehaviorSubject<Navigation>(emptyNavigation);

  navigation = emptyNavigation;

  private navigationSubscription: Subscription;

  constructor(
    private iconService: IconService,
    private navigationService: NavigationService,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.navigationSubscription = this.navigationService.current.subscribe(
      navigation => {
        if (this.navigation !== navigation) {
          this.navigation = navigation;
          this.cd.markForCheck();
        }
      }
    );
  }

  ngOnDestroy(): void {
    this.navigationSubscription.unsubscribe();
  }

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return this.iconService.load(item);
  }
}
