// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  OnDestroy,
  OnInit,
  HostListener,
} from '@angular/core';
import { Subscription } from 'rxjs';
import { Navigation, NavigationChild } from '../../../models/navigation';
import { IconService } from '../../../../shared/services/icon/icon.service';
import {
  Module,
  NavigationService,
  Selection,
} from '../../../../shared/services/navigation/navigation.service';
import { Router } from '@angular/router';
import { ThemeService } from '../../../../shared/services/theme/theme.service';

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
  collapsed = false;
  showLabels = true;
  selectedItem: Selection = { module: 0, index: -1 };
  flyoutIndex = -1;
  navigation = emptyNavigation;
  modules: Module[] = [];
  currentModule: Module;

  private subscriptionModules: Subscription;
  private subscriptionSelectedItem: Subscription;
  private subscriptionCollapsed: Subscription;
  private subscriptionShowLabels: Subscription;

  constructor(
    private iconService: IconService,
    private navigationService: NavigationService,
    private router: Router,
    private themeService: ThemeService,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.subscriptionModules = this.navigationService.modules.subscribe(
      modules => {
        this.modules = modules;
        this.currentModule = this.modules[this.selectedItem.module];
        this.cd.markForCheck();
      }
    );

    this.subscriptionSelectedItem = this.navigationService.selectedItem.subscribe(
      selection => {
        if (
          this.selectedItem.index !== selection.index ||
          this.selectedItem.module !== selection.module
        ) {
          this.selectedItem = selection;
          this.currentModule = this.modules[this.selectedItem.module];
          this.cd.markForCheck();
        }
      }
    );

    this.subscriptionCollapsed = this.navigationService.collapsed.subscribe(
      col => {
        if (this.collapsed !== col) {
          this.collapsed = col;
          this.cd.markForCheck();
        }
      }
    );

    this.subscriptionShowLabels = this.navigationService.showLabels.subscribe(
      col => {
        if (this.showLabels !== col) {
          this.showLabels = col;
          this.cd.markForCheck();
        }
      }
    );
  }

  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (event.key === 'T' && event.ctrlKey) {
      event.preventDefault();
      event.cancelBubble = true;
      this.themeService.switchTheme();
    } else if (event.key === 'b' && event.ctrlKey) {
      event.preventDefault();
      event.cancelBubble = true;
      this.updateNavCollapsed(!this.collapsed);
    } else if (event.key === 'L' && event.ctrlKey) {
      event.preventDefault();
      event.cancelBubble = true;
      this.navigationService.showLabels.next(!this.showLabels);
    }
  }

  identifyTab(index: number): string {
    return this.modules && this.modules.length > index
      ? this.modules[index].name
      : index.toString();
  }

  ngOnDestroy(): void {
    this.subscriptionModules?.unsubscribe();
    this.subscriptionSelectedItem?.unsubscribe();
    this.subscriptionCollapsed?.unsubscribe();
    this.subscriptionShowLabels?.unsubscribe();
  }

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  itemIcon(item: NavigationChild): string {
    return item.iconName;
  }

  formatPath(path: string): string {
    if (path && !path.startsWith('/')) {
      return '/' + path;
    }

    return path;
  }

  openPopup(index: number) {
    this.setNavState(true, index);
    this.setLastSelection(index);
  }

  closePopups(index) {
    this.flyoutIndex = -1;
    this.setLastSelection(index);
  }

  setNavState($event, state: number) {
    if ($event) {
      this.setLastSelection(state);
    }
  }

  shouldExpand(index: number) {
    if (this.collapsed) {
      return index === this.flyoutIndex;
    } else return index === this.selectedItem.index;
  }

  updateNavCollapsed(value: boolean): void {
    this.navigationService.collapsed.next(value);
  }

  setLastSelection(index) {
    if (this.selectedItem.index !== index) {
      this.navigationService.selectedItem.next({
        module: this.selectedItem.module,
        index,
      });
    }
  }

  setModule(module: number): void {
    this.navigationService.selectedItem.next({
      module,
      index: this.selectedItem.index,
    });
    if (this.modules[module].path) {
      this.router.navigateByUrl(this.modules[module].path);
    }
  }
}
