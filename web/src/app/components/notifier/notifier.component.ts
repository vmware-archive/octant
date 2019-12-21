// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  NotifierService,
  NotifierSignalType,
} from 'src/app/services/notifier/notifier.service';
import { Subscription } from 'rxjs';
import _ from 'lodash';

@Component({
  selector: 'app-notifier',
  templateUrl: './notifier.component.html',
  styleUrls: ['./notifier.component.scss'],
})
export class NotifierComponent implements OnInit, OnDestroy {
  private signalSubscription: Subscription;
  loading = false;
  error: string;
  warning: string;
  info: string;

  constructor(private notifierService: NotifierService) {}

  ngOnInit() {
    this.signalSubscription = this.notifierService.globalSignalsStream
      .pipe(untilDestroyed(this))
      .subscribe(currentSignals => {
        const lastLoadingSignal = _.findLast(currentSignals, {
          type: NotifierSignalType.LOADING,
        });
        this.loading = !!lastLoadingSignal;

        const lastWarningSignal = _.findLast(currentSignals, {
          type: NotifierSignalType.WARNING,
        });
        this.warning = lastWarningSignal
          ? (lastWarningSignal.data as string)
          : '';

        const lastErrorSignal = _.findLast(currentSignals, {
          type: NotifierSignalType.ERROR,
        });
        this.error = lastErrorSignal ? (lastErrorSignal.data as string) : '';

        const lastInfoSignal = _.findLast(currentSignals, {
          type: NotifierSignalType.INFO,
        });
        this.info = lastInfoSignal ? (lastInfoSignal.data as string) : '';
      });
  }

  ngOnDestroy(): void {
    if (this.signalSubscription) {
      this.signalSubscription.unsubscribe();
    }
  }
}
