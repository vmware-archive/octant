// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  Component,
  Input,
  OnChanges,
  OnInit,
  SimpleChanges,
} from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import {
  View,
  ButtonGroupView,
  PathItem,
} from 'src/app/modules/shared/models/content';
import { ViewService } from '../../../services/view/view.service';
import { WebsocketService } from '../../../services/websocket/websocket.service';

interface Tab {
  name: string;
  view: View;
  accessor: string;
  isClosable?: boolean;
}

@Component({
  selector: 'app-object-tabs',
  templateUrl: './tabs.component.html',
  styleUrls: ['./tabs.component.scss'],
})
export class TabsComponent implements OnChanges, OnInit {
  @Input() title: PathItem[];
  @Input() views: View[];
  @Input() payloads: [{ [key: string]: string }];
  @Input() closable: boolean;
  @Input() extView: boolean;
  @Input() buttonGroup: ButtonGroupView;

  tabs: Tab[] = [];
  activeTab: string;
  activeTabIndex: number;
  closingTab: boolean;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private viewService: ViewService,
    private wss: WebsocketService
  ) {}

  ngOnInit() {
    const { fragment } = this.activatedRoute.snapshot;
    if (fragment) {
      this.activeTab = fragment;
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.views.currentValue) {
      const views = changes.views.currentValue as View[];
      this.tabs = views.map(view => {
        const title = this.viewService.viewTitleAsText(view);
        return {
          name: title,
          view,
          accessor: view.metadata.accessor,
        };
      });

      if (this.extView && this.tabs.length > 0) {
        // Initial load if there are existing tabs
        if (this.activeTabIndex === null) {
          this.activeTabIndex = this.tabs.length - 1;
        }

        if (!changes.views.isFirstChange()) {
          const preViews = changes.views.previousValue as View[];
          // Focus new tab
          if (views.length > preViews.length && !this.closingTab) {
            this.activeTabIndex = this.tabs.length - 1;
          }
          this.closingTab = false;
        }
        this.activeTab = this.tabs[this.activeTabIndex].accessor;
        this.setMarker(this.activeTab);
      }
    }
  }

  identifyTab(index: number, item: Tab): string {
    return item.name;
  }

  clickTab(tabAccessor: string) {
    if (this.activeTab === name) {
      return;
    }
    this.activeTab = tabAccessor;
    this.setMarker(tabAccessor);
  }

  closeTab(tabAccessor: string) {
    const tabIndex = this.tabs.findIndex(tab => tab.accessor === tabAccessor);
    if (tabIndex > -1) {
      if (this.payloads[tabIndex]) {
        const payload = this.payloads[tabIndex];
        this.wss.sendMessage('action.octant.dev/performAction', payload);
      }

      this.tabs = [
        ...this.tabs.slice(0, tabIndex),
        ...this.tabs.slice(tabIndex + 1),
      ];

      this.closingTab = true;

      switch (this.tabs.length > 0) {
        // Closed active tab
        case tabIndex === this.activeTabIndex:
          // Closed right most active tab
          if (tabIndex === this.tabs.length) {
            this.activeTabIndex -= 1;
          } else {
            this.activeTabIndex = tabIndex;
          }
          break;
        // Closed left of active tab
        case tabIndex < this.activeTabIndex:
          this.activeTabIndex -= 1;
          break;
        // Closed right of active tab
        case tabIndex > this.activeTabIndex:
          // no-op
          break;
        default:
        // no-op
      }
    }
    this.activeTab = this.tabs[this.activeTabIndex].accessor;
    this.setMarker(this.activeTab);
  }

  private setMarker(tabAccessor: string) {
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      replaceUrl: true,
      fragment: tabAccessor,
    });
  }
}
