// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import findLast from 'lodash/findLast';
import { Subscription } from 'rxjs';
import { Alert } from 'src/app/modules/shared/models/content';
import {
  NotifierService,
  NotifierSignalType,
} from 'src/app/modules/shared/notifier/notifier.service';

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
  success: string;
  alertConfig: Alert;

  constructor(private notifierService: NotifierService) {}

  ngOnInit() {
    this.signalSubscription =
      this.notifierService.globalSignalsStream.subscribe(currentSignals => {
        const lastLoadingSignal = findLast(currentSignals, {
          type: NotifierSignalType.LOADING,
        });
        this.loading = !!lastLoadingSignal;

        const lastWarningSignal = findLast(currentSignals, {
          type: NotifierSignalType.WARNING,
        });
        this.warning = lastWarningSignal
          ? (lastWarningSignal.data as string)
          : '';

        const lastErrorSignal = findLast(currentSignals, {
          type: NotifierSignalType.ERROR,
        });
        this.error = lastErrorSignal ? (lastErrorSignal.data as string) : '';

        const lastInfoSignal = findLast(currentSignals, {
          type: NotifierSignalType.INFO,
        });
        this.info = lastInfoSignal ? (lastInfoSignal.data as string) : '';

        const lastSuccessSignal = findLast(currentSignals, {
          type: NotifierSignalType.SUCCESS,
        });
        this.success = lastSuccessSignal
          ? (lastSuccessSignal.data as string)
          : '';

        this.setAlert();
      });
  }

  ngOnDestroy(): void {
    if (this.signalSubscription) {
      this.signalSubscription.unsubscribe();
    }
  }

  hasAlert(): boolean {
    return !!(this.info || this.success || this.error || this.warning);
  }

  setAlert(): void {
    let status: string, message: string;
    let closable: boolean;

    if (this.warning) {
      status = 'warning';
      message = this.warning;
    } else if (this.error) {
      status = 'error';
      message = this.error;
    } else if (this.success) {
      status = 'success';
      message = this.success;
    } else {
      status = 'info';
      message = this.info;
      closable = true;
    }

    this.alertConfig = {
      status: status,
      type: 'banner',
      message: message,
      closable: !!closable,
    };
  }
}
