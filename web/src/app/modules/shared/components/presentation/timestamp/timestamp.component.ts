// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  Input,
  OnChanges,
  OnDestroy,
  SimpleChanges,
  ChangeDetectionStrategy,
  OnInit,
  ChangeDetectorRef,
} from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView, View } from 'src/app/modules/shared/models/content';
import { Subscription } from 'rxjs';
import { ContentService } from '../../../services/content/content.service';

@Component({
  selector: 'app-view-timestamp',
  templateUrl: './timestamp.component.html',
  styleUrls: ['./timestamp.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimestampComponent implements OnInit, OnChanges, OnDestroy {
  private v: TimestampView;

  @Input() set view(v: View) {
    this.v = v as TimestampView;
  }
  get view() {
    return this.v;
  }

  timestamp: number;
  humanReadable: string;
  age: string;
  scrollPosition = 0;
  private contentSubscription: Subscription;

  constructor(
    private contentService: ContentService,
    private cd: ChangeDetectorRef
  ) {
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

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as TimestampView;

      this.timestamp = view.config.timestamp;
      this.humanReadable =
        dayjs(this.timestamp * 1000)
          .utc()
          .format('llll') + ' UTC';
    }
  }

  ngOnDestroy(): void {
    this.timestamp = null;
    if (this.contentSubscription) {
      this.contentSubscription.unsubscribe();
    }
  }
}
