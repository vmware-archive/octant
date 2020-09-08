// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { AfterViewChecked, Component, ElementRef, Input } from '@angular/core';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

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
export class OverflowSelectorsComponent implements AfterViewChecked {
  @Input() set selectors(selectors: Selector[]) {
    this.selectorsList = selectors;
    this.updateSelectors();
  }

  get selectors(): Selector[] {
    return this.selectorsList;
  }

  constructor(private rootElement: ElementRef) {}
  @Input() numberShownSelectors = 2;

  private selectorsList: Selector[];
  showSelectors: Selector[];
  overflowSelectors: Selector[];
  trackByIdentity = trackByIdentity;
  componentWidth = 0;

  private updateSelectors() {
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

  ngAfterViewChecked(): void {
    if (this.componentWidth !== this.rootElement.nativeElement.clientWidth) {
      this.numberShownSelectors =
        this.rootElement.nativeElement.clientWidth > 150 ? 2 : 1;
      this.updateSelectors();
      this.componentWidth = this.rootElement.nativeElement.clientWidth;
    }
  }
}
