// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  ElementRef,
  HostListener,
  OnDestroy,
  OnInit,
} from '@angular/core';
import '@cds/core/modal/register.js';
import { Router } from '@angular/router';
import { BehaviorSubject, Subject, Subscription } from 'rxjs';
import { Navigation, NavigationChild } from '../../../models/navigation';
import { NavigationService } from '../../../../shared/services/navigation/navigation.service';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { NamespaceService } from 'src/app/modules/shared/services/namespace/namespace.service';

const emptyNavigation: Navigation = {
  sections: [],
  defaultPath: '',
};

interface Destination {
  title: string;
  type: string;
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
  searchingNamespace = false;

  destinations: Destination[];
  namespaceDestinations: Destination[];
  filteredDestinations: Destination[];
  currentDestination = '';

  helperText = `Search namespaces by starting with `;

  input = '';
  inputChanged: Subject<string> = new Subject<string>();

  activeIndex = 0;
  styledShadowDom = false;

  private navigationSubscription: Subscription;
  private namespaceSubscription: Subscription;

  constructor(
    private navigationService: NavigationService,
    private namespaceService: NamespaceService,
    private router: Router,
    private el: ElementRef
  ) {
    // wait a bit before reacting to user input
    this.inputChanged
      .pipe(debounceTime(150), distinctUntilChanged())
      .subscribe(f => this.updateFilteredDestinations(f));
  }

  ngOnInit() {
    this.namespaceSubscription = this.namespaceService.availableNamespaces.subscribe(
      namespaces => {
        this.namespaceDestinations = this.buildNamespaceDestinations(
          namespaces
        );
      }
    );
    this.navigationSubscription = this.navigationService.current.subscribe(
      navigation => {
        this.navigation = navigation;
        this.destinations = this.buildDestinations(navigation);
      }
    );
  }

  ngOnDestroy(): void {
    if (this.navigationSubscription) {
      this.navigationSubscription.unsubscribe();
    }
    if (this.namespaceSubscription) {
      this.namespaceSubscription.unsubscribe();
    }
  }

  identifyNavigationItem(index: number, item: NavigationChild): string {
    return item.title;
  }

  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (event.key === 'Enter' && event.ctrlKey) {
      this.resetModal();
      this.toggleQuickSwitcher();

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

  private buildNamespaceDestinations(namespaces: string[]): Destination[] {
    const nsDestinations = [];
    namespaces.forEach(namespace => {
      nsDestinations.push({
        title: namespace,
        path: 'overview/namespace/' + namespace,
        keywords: [namespace],
      });
    });
    return nsDestinations;
  }

  private recBuildDestinations(titleAcc: string, keywordAcc, item) {
    if (titleAcc !== '') {
      item.type = titleAcc;
    }
    if (!item.children) {
      let k = keywordAcc.slice(0);
      k.push(item.title);
      k = k.flatMap(i => i.split(' '));
      return [
        {
          title: item.title,
          type: item.type,
          path: item.path,
          keywords: k,
          active: false,
        },
      ];
    }
    keywordAcc.push(item.title);
    return item.children.flatMap(child =>
      this.recBuildDestinations(item.title, keywordAcc, child)
    );
  }

  identifyDestinationItem(_: number, item: Destination): string {
    return item.title;
  }

  onInputChange(input: string) {
    this.inputChanged.next(input);
  }

  onEnter() {
    const d = this.filteredDestinations[this.activeIndex].path;
    this.router.navigateByUrl(d);
    this.toggleQuickSwitcher();
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
      this.searchingNamespace = false;
      return;
    }

    if (this.input.startsWith('!')) {
      this.searchingNamespace = true;
      filter = filter.substring(1);
      this.filteredDestinations = this.namespaceDestinations.filter(d => {
        return (
          d.title !== this.namespaceService.activeNamespace.value &&
          d.title.toLowerCase().includes(filter.toLowerCase())
        );
      });
      return;
    }
    this.searchingNamespace = false;
    this.filteredDestinations = this.destinations.filter(d => {
      const lk = d.keywords.map(k => k.toLowerCase());
      return lk.findIndex(k => k.includes(filter.toLowerCase())) !== -1;
    });
  }

  private resetModal() {
    this.input = '';
    this.filteredDestinations = this.destinations;
    this.activeIndex = 0;
    this.searchingNamespace = false;
  }

  toggleQuickSwitcher(): void {
    const qcModal = document.getElementById('quick-switcher-modal');
    qcModal.hidden = !qcModal.hidden;

    // Add styling to prevent modal from moving as number of results update
    if (!this.styledShadowDom) {
      const style = document.createElement('style');
      style.innerHTML =
        '.modal-dialog { position: fixed !important; top: 4rem; }';
      qcModal.shadowRoot.appendChild(style);
      this.styledShadowDom = true;
    }
  }
}
