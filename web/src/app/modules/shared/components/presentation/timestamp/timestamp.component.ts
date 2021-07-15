// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy } from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-timestamp',
  templateUrl: './timestamp.component.html',
  styleUrls: ['./timestamp.component.scss'],
})
export class TimestampComponent
  extends AbstractViewComponent<TimestampView>
  implements OnDestroy
{
  timestamp: number;
  humanReadable: string;

  constructor() {
    super();
    dayjs.extend(utc);
    dayjs.extend(LocalizedFormat);
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
  }
}
