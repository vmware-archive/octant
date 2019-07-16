// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
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

  constructor(private notifierService: NotifierService) {}

  ngOnInit() {
    this.signalSubscription = this.notifierService.globalSignalsStream.subscribe(
      currentSignals => {
        const lastLoadingSignal = _.findLast(currentSignals, {
          type: NotifierSignalType.LOADING,
        });
        this.loading = lastLoadingSignal ? true : false;

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
      }
    );
  }

  onWarningClose() {
    this.warning = '';
    // TODO: remove warning from signals queue?
  }

  ngOnDestroy(): void {
    if (this.signalSubscription) {
      this.signalSubscription.unsubscribe();
    }
  }
}
