// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed } from '@angular/core/testing';
import {
  NotifierService,
  NotifierSignal,
  NotifierSignalType,
} from './notifier.service';

describe('NotifierService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: NotifierService = TestBed.inject(NotifierService);
    expect(service).toBeTruthy();
    expect(service.globalSignalsStream.getValue()).toEqual([]);
  });

  it('should be able to add signals', () => {
    const service = new NotifierService();
    const loadingSignalID = service.pushSignal(
      NotifierSignalType.LOADING,
      true
    );
    const errorSignalID = service.pushSignal(
      NotifierSignalType.ERROR,
      'You are doing it wrong'
    );
    const observedSignals = service.globalSignalsStream.getValue();
    const expectedSignals: Array<NotifierSignal> = [
      {
        id: loadingSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.LOADING,
        data: true,
      },
      {
        id: errorSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.ERROR,
        data: 'You are doing it wrong',
      },
    ];
    expect(observedSignals).toEqual(expectedSignals);
  });

  it('should be able to remove signals', () => {
    const service = new NotifierService();

    const warningSignalID = service.pushSignal(
      NotifierSignalType.WARNING,
      'nope nope nope'
    );
    const loadingSignalID = service.pushSignal(
      NotifierSignalType.LOADING,
      false
    );
    const errorSignalID = service.pushSignal(
      NotifierSignalType.ERROR,
      'No more worky'
    );

    const initialObservedSignals = service.globalSignalsStream.getValue();
    const initialExpectedSignals: Array<NotifierSignal> = [
      {
        id: warningSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.WARNING,
        data: 'nope nope nope',
      },
      {
        id: loadingSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.LOADING,
        data: false,
      },
      {
        id: errorSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.ERROR,
        data: 'No more worky',
      },
    ];
    expect(initialObservedSignals).toEqual(initialExpectedSignals);

    service.removeSignals([loadingSignalID, '', null]);

    const currentObservedSignals = service.globalSignalsStream.getValue();
    const currentExpectedSignals: Array<NotifierSignal> = [
      {
        id: warningSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.WARNING,
        data: 'nope nope nope',
      },
      {
        id: errorSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.ERROR,
        data: 'No more worky',
      },
    ];
    expect(currentObservedSignals).toEqual(currentExpectedSignals);
  });

  it('should be able to create a signal session', () => {
    const service = new NotifierService();

    const warningSignalID = service.pushSignal(
      NotifierSignalType.WARNING,
      'nope nope nope'
    );

    const session = service.createSession();

    const loadingSignalID = session.pushSignal(
      NotifierSignalType.LOADING,
      false
    );
    const errorSignalID = session.pushSignal(
      NotifierSignalType.ERROR,
      'No more worky'
    );

    const initialObservedSignals = service.globalSignalsStream.getValue();
    const initialExpectedSignals: Array<NotifierSignal> = [
      {
        id: warningSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.WARNING,
        data: 'nope nope nope',
      },
      {
        id: loadingSignalID,
        sessionID: session.id,
        type: NotifierSignalType.LOADING,
        data: false,
      },
      {
        id: errorSignalID,
        sessionID: session.id,
        type: NotifierSignalType.ERROR,
        data: 'No more worky',
      },
    ];
    expect(initialObservedSignals).toEqual(initialExpectedSignals);

    session.removeAllSignals();

    const currentObservedSignals = service.globalSignalsStream.getValue();
    const currentExpectedSignals: Array<NotifierSignal> = [
      {
        id: warningSignalID,
        sessionID: 'baseSignal',
        type: NotifierSignalType.WARNING,
        data: 'nope nope nope',
      },
    ];
    expect(currentObservedSignals).toEqual(currentExpectedSignals);
  });
});
