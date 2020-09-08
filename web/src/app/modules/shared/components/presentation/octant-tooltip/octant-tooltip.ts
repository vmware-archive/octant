// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectorRef,
  Component,
  Input,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { Subscription } from 'rxjs';
import { ContentService } from '../../../services/content/content.service';

@Component({
  selector: 'app-octant-tooltip',
  templateUrl: './octant-tooltip.html',
  styleUrls: ['./octant-tooltip.scss'],
})
export class OctantTooltipComponent implements OnInit, OnDestroy {
  scrollPosition = 0;
  private contentSubscription: Subscription;

  @Input()
  tooltipText: string;

  constructor(
    private contentService: ContentService,
    private cd: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.contentSubscription = this.contentService.viewScrollPos.subscribe(
      position => {
        this.scrollPosition = position;
        this.cd.markForCheck();
      }
    );
  }

  getScrollPos() {
    return `${-this.scrollPosition - 64}px`;
  }

  ngOnDestroy(): void {
    if (this.contentSubscription) {
      this.contentSubscription.unsubscribe();
    }
  }
}
