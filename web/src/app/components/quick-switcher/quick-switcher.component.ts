// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  OnInit,
  OnDestroy,
  HostListener,
  ElementRef,
} from '@angular/core';
import { Router } from '@angular/router';
import { Subject, BehaviorSubject } from 'rxjs';
import { Navigation, NavigationChild } from '../../models/navigation';
import { NavigationService } from '../../modules/overview/services/navigation/navigation.service';
import { untilDestroyed } from 'ngx-take-until-destroy';
import { distinctUntilChanged, debounceTime } from 'rxjs/operators';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

interface Destination {
  title: string;
  path: string;
  keywords: string[];
}

@Component({
  selector: 'app-quick-switcher',
  templateUrl: './quick-switcher.component.html',
  styleUrls: ['./quick-switcher.component.scss'],
})
export class QuickSwitcherComponent implements OnInit, OnDestroy {
  behavior = new BehaviorSubject<Navigation>(emptyNavigation);

  navigation: Navigation = emptyNavigation;

  opened = false;

  destinations: Destination[];
  filteredDestinations: Destination[];
  currentDestination = '';

  input = '';
  inputChanged: Subject<string> = new Subject<string>();

  activeIndex = 0;

  constructor(
    private navigationService: NavigationService,
    private router: Router,
    private el: ElementRef
  ) {
    // wait a bit before reacting to user input
    this.inputChanged
      .pipe(
        debounceTime(150),
        distinctUntilChanged()
      )
      .subscribe(f => this.updateFilteredDestinations(f));
  }

  ngOnInit() {
    this.navigationService.current
      .pipe(untilDestroyed(this))
      .subscribe(navigation => {
        this.navigation = navigation;
        this.destinations = this.buildDestinations(navigation);
      });
  }

  ngOnDestroy() {}

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (this.opened) {
      return;
    }
    if (event.key === 'k' && event.ctrlKey) {
      this.resetModal();
      this.opened = true;

      // TODO(abrand): Hack for focusing on input. Need to figure this out. (GH#508)
      const el = this.el;
      setTimeout(() => {
        el.nativeElement.querySelector('.filter-input').focus();
      }, 250);
    }
  }

  // buildDestinations recursively builds an array of destinations based on the
  // information obtained from the navigation service
  private buildDestinations(navigation) {
    return navigation.sections.flatMap(section =>
      this.recBuildDestinations('', [], section)
    );
  }

  private recBuildDestinations(titleAcc: string, keywordAcc, item) {
    let title = item.title;
    if (titleAcc !== '') {
      title = titleAcc + ' -> ' + item.title;
    }
    if (!item.children) {
      let k = keywordAcc.slice(0);
      k.push(item.title);
      k = k.flatMap(i => i.split(' '));
      return [{ title, path: item.path, keywords: k, active: false }];
    }
    keywordAcc.push(item.title);
    return item.children.flatMap(child =>
      this.recBuildDestinations(title, keywordAcc, child)
    );
  }

  private handleEvent = (message: MessageEvent) => {
    const data = JSON.parse(message.data);
    this.behavior.next(data);
  };

  identifyDestinationItem(_: number, item: Destination): string {
    return item.title;
  }

  navigateTo(destination: string) {
    console.log(destination);
  }

  onInputChange(input: string) {
    this.inputChanged.next(input);
  }

  onEnter() {
    const d = this.filteredDestinations[this.activeIndex].path;
    this.router.navigateByUrl(d);
    this.opened = false;
    this.resetModal();
  }

  onKeyUp(event: KeyboardEvent) {
    if (event.key === 'ArrowDown') {
      event.preventDefault();
      this.activeIndex = Math.min(
        this.activeIndex + 1,
        this.filteredDestinations.length - 1
      );
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      this.activeIndex = Math.max(this.activeIndex - 1, 0);
    }
  }

  onKeyDown(event: KeyboardEvent) {
    if (event.key === 'ArrowUp' || event.key === 'ArrowDown') {
      event.preventDefault();
    }
  }

  updateFilteredDestinations(filter: string) {
    this.activeIndex = 0;
    if (filter === '') {
      this.filteredDestinations = this.destinations;
      return;
    }
    this.filteredDestinations = this.destinations.filter(d => {
      const lk = d.keywords.map(k => k.toLowerCase());
      return lk.findIndex(k => k.includes(filter)) !== -1;
    });
  }

  private resetModal() {
    this.input = '';
    this.filteredDestinations = this.destinations;
    this.activeIndex = 0;
  }
}
