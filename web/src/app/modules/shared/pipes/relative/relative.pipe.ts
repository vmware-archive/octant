/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import {
  ChangeDetectorRef,
  NgZone,
  OnDestroy,
  Pipe,
  PipeTransform,
} from '@angular/core';

const changeDetectionFrequency = (seconds: number) => {
  switch (true) {
    case seconds < 60:
      return 1;
    case seconds < 3600:
      return 60;
    case seconds < 86400:
      return 600;
    default:
      return 3600;
  }
};

@Pipe({
  name: 'relative',
  pure: false,
})

/**
 * RelativePipe converts a timestamp to a relative time from the current time.
 * If no date is supplied, it will use the current date.
 *
 * @param ts timestamp in seconds since epoch
 * @param base optional date to calculate from
 */
export class RelativePipe implements PipeTransform, OnDestroy {
  private timer: number;

  constructor(
    private changeDetectorRef: ChangeDetectorRef,
    private ngZone: NgZone
  ) {}

  transform(ts: number, base?: Date): string {
    this.removeTimer();

    let now: Date;
    if (base) {
      now = base;
    } else {
      now = new Date();
    }

    const then = now.getTime() / 1000 - ts;

    const updateInterval = changeDetectionFrequency(then) * 1000;

    this.timer = this.ngZone.runOutsideAngular(() => {
      if (typeof window !== 'undefined') {
        return window.setTimeout(() => {
          this.ngZone.run(() => {
            if (this.changeDetectorRef) {
              this.changeDetectorRef.markForCheck();
            }
          });
        }, updateInterval);
      }
      return null;
    });

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

  ngOnDestroy(): void {
    this.removeTimer();
  }

  private removeTimer() {
    if (this.timer) {
      window.clearTimeout(this.timer);
      this.timer = null;
    }
  }
}
