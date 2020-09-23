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
  navExpandedState: any;
  selectedItem: Selection = { module: 0, index: -1 };
  flyoutIndex = -1;
  navigation = emptyNavigation;
  modules: Module[] = [];

  private navigationSubscription: Subscription;

  constructor(
    private iconService: IconService,
    private navigationService: NavigationService,
    private router: Router,
    private themeService: ThemeService,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.navigationSubscription = this.navigationService.modules.subscribe(
      modules => {
        this.modules = modules;
        this.cd.markForCheck();
      }
    );

    this.navigationSubscription = this.navigationService.selectedItem.subscribe(
      selection => {
        if (
          this.selectedItem.index !== selection.index ||
          this.selectedItem.module !== selection.module
        ) {
          this.selectedItem = selection;
          this.cd.markForCheck();
        }
      }
    );

    this.navigationSubscription = this.navigationService.expandedState.subscribe(
      state => {
        if (this.navExpandedState !== state) {
          this.navExpandedState = state;
          this.cd.markForCheck();
        }
      }
    );

    this.navigationSubscription = this.navigationService.collapsed.subscribe(
      col => {
        if (this.collapsed !== col) {
          this.collapsed = col;
          this.cd.markForCheck();
        }
      }
    );

    this.navigationSubscription = this.navigationService.showLabels.subscribe(
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
      this.themeService.switchTheme();
    } else if (event.key === 'N' && event.ctrlKey) {
      this.updateNavCollapsed(!this.collapsed);
    } else if (event.key === 'L' && event.ctrlKey) {
      this.navigationService.showLabels.next(!this.showLabels);
    }
  }

  identifyTab(index: number): string {
    return this.modules && this.modules.length > index
      ? this.modules[index].name
      : index.toString();
  }

  getSelectedSections() {
    const currentModule = this.modules[this.selectedItem.module];
    return currentModule ? currentModule.children : [];
  }

  getModuleTitle() {
    const currentModule = this.modules[this.selectedItem.module];
    return currentModule ? currentModule.title + ' Module' : '';
  }

  getModuleDescription() {
    const currentModule = this.modules[this.selectedItem.module];
    return currentModule ? currentModule.description : '';
  }

  ngOnDestroy(): void {
    if (this.navigationSubscription) {
      this.navigationSubscription.unsubscribe();
    }
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
    this.clearExpandedState();
    this.setNavState(true, index);
    this.setLastSelection(index);
  }

  closePopups(index) {
    this.clearExpandedState();
    this.flyoutIndex = -1;
    this.setLastSelection(index);
  }

  setExpandedState(index, state) {
    this.navExpandedState[index] = state;
    this.navigationService.expandedState.next(this.navExpandedState);
  }

  clearExpandedState() {
    this.navExpandedState = {};
    this.navigationService.expandedState.next(this.navExpandedState);
  }

  setNavState($event, state: number) {
    if (this.collapsed) {
      this.setLastSelection(state);
    } else {
      this.setExpandedState(state, $event);
      if ($event && this.selectedItem.index !== state) {
        // collapse previously selected group
        if (this.selectedItem) {
          this.setExpandedState(this.selectedItem.index, false);
        }
        this.setLastSelection(state);
      }
    }
  }

  shouldExpand(index: number) {
    if (this.collapsed) {
      return index === this.flyoutIndex;
    } else if (index.toString() in this.navExpandedState) {
      return this.navExpandedState[index];
    }
    return false;
  }

  updateNavCollapsed(value: boolean): void {
    this.collapsed = value;
    this.navigationService.collapsed.next(value);
    this.setExpandedState(this.selectedItem.index, false);
  }

  setLastSelection(index) {
    this.selectedItem.index = index;
    this.navigationService.selectedItem.next({
      module: this.selectedItem.module,
      index,
    });
  }

  setModule(module: number): void {
    if (this.selectedItem) {
      this.setExpandedState(this.selectedItem.index, false);
    }

    this.selectedItem.module = module;
    this.navigationService.selectedItem.next({
      module,
      index: this.selectedItem.index,
    });
    if (this.modules[module].path) {
      this.router.navigateByUrl(this.modules[module].path);
    }
  }
}
