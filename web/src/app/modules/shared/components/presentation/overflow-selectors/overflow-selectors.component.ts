// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { ContentService } from '../../../services/content/content.service';
import { Subscription } from 'rxjs';

interface Selector {
  metadata: {
    type: string;
  };
  config: {
    key: string;
    value: string;
  };
}

@Component({
  selector: 'app-overflow-selectors',
  templateUrl: './overflow-selectors.component.html',
  styleUrls: ['./overflow-selectors.component.scss'],
})
export class OverflowSelectorsComponent implements OnInit, OnDestroy {
  @Input() numberShownSelectors = 2;
  @Input() set selectors(selectors: Selector[]) {
    this.selectorsList = selectors;

    if (this.numberShownSelectors <= this.selectorsList.length) {
      this.showSelectors = this.selectorsList.slice(
        0,
        this.numberShownSelectors
      );
      this.overflowSelectors = this.selectorsList.slice(
        this.numberShownSelectors
      );
    } else {
      this.showSelectors = this.selectorsList;
    }
  }
  get selectors(): Selector[] {
    return this.selectorsList;
  }

  private selectorsList: Selector[];
  showSelectors: Selector[];
  overflowSelectors: Selector[];
  trackByIdentity = trackByIdentity;
  scrollPosition = 0;
  private contentSubscription: Subscription;

  constructor(private contentService: ContentService) {}

  ngOnInit() {
    this.contentSubscription = this.contentService.viewScrollPos.subscribe(
      position => {
        this.scrollPosition = position;
        console.log('POSITION ->>', position);
      }
    );
  }

  ngOnDestroy() {
    this.contentSubscription.unsubscribe();
  }

  getScrollPos() {
    return `${-this.scrollPosition - 64}px`;
  }
}
