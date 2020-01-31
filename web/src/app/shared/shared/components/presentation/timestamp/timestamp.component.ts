// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView, View } from 'src/app/models/content';

@Component({
  selector: 'app-view-timestamp',
  templateUrl: './timestamp.component.html',
  styleUrls: ['./timestamp.component.scss'],
})
export class TimestampComponent implements OnChanges {
  private v: TimestampView;

  @Input() set view(v: View) {
    this.v = v as TimestampView;
  }
  get view() {
    return this.v;
  }

  humanReadable: string;
  age: string;

  constructor() {
    dayjs.extend(utc);
    dayjs.extend(LocalizedFormat);
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as TimestampView;

      const timestamp = view.config.timestamp;
      this.humanReadable =
        dayjs(timestamp * 1000)
          .utc()
          .format('LLLL') + ' UTC';
      this.age = this.summarizeTimestamp(timestamp);
    }
  }

  /**
   * summarizeTimestamp converts a timestamp to a relative time from the current time.
   * If no date is supplied, it will use the current date.
   *
   * @param ts timestamp in seconds since epoch
   * @param base optional date to calculate from
   */
  summarizeTimestamp(ts: number, base?: Date): string {
    let now: Date;
    if (base) {
      now = base;
    } else {
      now = new Date();
    }

    const then = now.getTime() / 1000 - ts;

    if (then > 86400) {
      return `${Math.floor(then / 86400)}d`;
    } else if (then > 3600) {
      return `${Math.floor(then / 3600)}h`;
    } else if (then > 60) {
      return `${Math.floor(then / 60)}m`;
    } else {
      return `${Math.floor(then)}s`;
    }
  }
}
