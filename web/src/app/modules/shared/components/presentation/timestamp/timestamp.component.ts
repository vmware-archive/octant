// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView } from 'src/app/modules/shared/models/content';
import { Subscription } from 'rxjs';
import { ContentService } from '../../../services/content/content.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-timestamp',
  templateUrl: './timestamp.component.html',
  styleUrls: ['./timestamp.component.scss'],
})
export class TimestampComponent
  extends AbstractViewComponent<TimestampView>
  implements OnInit, OnDestroy {
  timestamp: number;
  humanReadable: string;
  age: string;
  scrollPosition = 0;
  private contentSubscription: Subscription;

  constructor(
    private contentService: ContentService,
    private cd: ChangeDetectorRef
  ) {
    super();
    dayjs.extend(utc);
    dayjs.extend(LocalizedFormat);
  }

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

  update() {
    this.timestamp = this.v.config.timestamp;
    this.humanReadable =
      dayjs(this.timestamp * 1000)
        .utc()
        .format('llll') + ' UTC';
  }

  ngOnDestroy(): void {
    this.timestamp = null;
    if (this.contentSubscription) {
      this.contentSubscription.unsubscribe();
    }
  }
}
