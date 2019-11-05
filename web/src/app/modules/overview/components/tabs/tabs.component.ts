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
import { View } from 'src/app/models/content';
import { ViewService } from '../../services/view/view.service';

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
  @Input() title: string;
  @Input() views: View[];
  @Input() iconName: string;

  tabs: Tab[] = [];
  activeTab: string;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private viewService: ViewService
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

      if (!this.activeTab) {
        this.activeTab = this.tabs[0].accessor;
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
      this.tabs = [
        ...this.tabs.slice(0, tabIndex),
        ...this.tabs.slice(tabIndex + 1),
      ];
    }
  }

  private setMarker(tabAccessor: string) {
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      replaceUrl: true,
      fragment: tabAccessor,
    });
  }
}
