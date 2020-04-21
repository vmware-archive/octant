// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  Input,
  OnChanges,
  OnDestroy,
  SimpleChanges,
  ViewChild,
  ChangeDetectionStrategy,
} from '@angular/core';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import LocalizedFormat from 'dayjs/plugin/localizedFormat';
import { TimestampView, View } from 'src/app/modules/shared/models/content';

@Component({
  selector: 'app-view-timestamp',
  templateUrl: './timestamp.component.html',
  styleUrls: ['./timestamp.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TimestampComponent implements OnChanges, OnDestroy {
  private v: TimestampView;

  @ViewChild('timestampRef', { static: true }) timestampRef: ElementRef;
  @Input() set view(v: View) {
    this.v = v as TimestampView;
  }
  get view() {
    return this.v;
  }

  timestamp: number;
  humanReadable: string;
  age: string;

  constructor() {
    dayjs.extend(utc);
    dayjs.extend(LocalizedFormat);
  }

  get tooltipPosition(): string {
    const gutterWidth = 300;
    const { nativeElement } = this.timestampRef;
    const timestampLeft = nativeElement.getBoundingClientRect().left;

    return timestampLeft > window.outerWidth - gutterWidth ? 'left' : 'right';
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
  }
}
