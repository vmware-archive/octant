// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import _ from 'lodash';

export enum NotifierSignalType {
  LOADING = 'LOADING',
  ERROR = 'ERROR',
  WARNING = 'WARNING',
  INFO = 'INFO',
}

export interface NotifierSignal {
  id: string;
  sessionID: string;
  type: NotifierSignalType;
  data: boolean | string;
}

export class NotifierSession {
  id: string;

  constructor(
    private globalSignalsStream: BehaviorSubject<NotifierSignal[]>,
    private uniqueIDPrefix: string
  ) {
    this.id = uniqueIDPrefix;
  }

  pushSignal(type: NotifierSignalType, data: boolean | string): string {
    const currentSignals = this.globalSignalsStream.getValue();
    const newSignalID = _.uniqueId(this.uniqueIDPrefix);
    const newSignal = {
      id: newSignalID,
      sessionID: this.uniqueIDPrefix,
      type,
      data,
    };
    this.globalSignalsStream.next([...currentSignals, newSignal]);
    return newSignalID;
  }

  removeSignal(id: string): boolean {
    const currentSignals = this.globalSignalsStream.getValue();
    const foundSignalIndex = _.findIndex(currentSignals, {
      id,
      sessionID: this.uniqueIDPrefix,
    });
    if (foundSignalIndex < 0) {
      return false;
    }

    const newSignalList = [...currentSignals];
    _.pullAt(newSignalList, foundSignalIndex);
    this.globalSignalsStream.next(newSignalList);
    return true;
  }

  removeSignals(ids: string[]): void {
    _.forEach(ids, (id: string) => {
      if (id) {
        this.removeSignal(id);
      }
    });
  }

  removeAllSignals(): void {
    const currentSignals = this.globalSignalsStream.getValue();
    const newSignalList = [...currentSignals];
    _.remove(newSignalList, { sessionID: this.uniqueIDPrefix });
    this.globalSignalsStream.next(newSignalList);
  }
}

@Injectable({
  providedIn: 'root',
})
export class NotifierService {
  baseSignalSession: NotifierSession;
  globalSignalsStream: BehaviorSubject<NotifierSignal[]> = new BehaviorSubject(
    []
  );

  constructor() {
    this.baseSignalSession = new NotifierSession(
      this.globalSignalsStream,
      'baseSignal'
    );
  }

  pushSignal(type: NotifierSignalType, data: boolean | string): string {
    return this.baseSignalSession.pushSignal(type, data);
  }

  removeSignal(id: string): boolean {
    return this.baseSignalSession.removeSignal(id);
  }

  removeSignals(ids: string[]): void {
    return this.baseSignalSession.removeSignals(ids);
  }

  createSession(): NotifierSession {
    return new NotifierSession(
      this.globalSignalsStream,
      _.uniqueId('signalSession')
    );
  }
}
