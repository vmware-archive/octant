// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView, View } from 'src/app/modules/shared/models/content';

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

  timestamp: number;
  humanReadable: string;

  constructor() {
    dayjs.extend(utc);
    dayjs.extend(LocalizedFormat);
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as TimestampView;

      this.timestamp = view.config.timestamp;
      this.humanReadable =
        dayjs(this.timestamp * 1000)
          .utc()
          .format('LLLL') + ' UTC';
    }
  }
}
