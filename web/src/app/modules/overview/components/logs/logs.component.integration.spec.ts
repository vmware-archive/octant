// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { HttpClient } from '@angular/common/http';
import { HttpTestingController } from '@angular/common/http/testing';
import { DebugElement } from '@angular/core';
import {
  async,
  ComponentFixture,
  discardPeriodicTasks,
  fakeAsync,
  TestBed,
} from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import _ from 'lodash';
import { LogEntry, LogsView } from 'src/app/models/content';
import getAPIBase from 'src/app/services/common/getAPIBase';
import { PodLogsService } from 'src/app/services/pod-logs/pod-logs.service';

import { LogsComponent } from './logs.component';

const API_BASE = getAPIBase();

function createTestLogsView(containers: string[]): LogsView {
  return {
    metadata: {
      type: 'logs',
      title: [],
      accessor: 'logs',
    },
    config: {
      namespace: 'default',
      name: 'cartpod',
      containers,
    },
  };
}

function createRandomLogEntry(): LogEntry {
  return {
    timestamp: '2019-05-06T18:59:06.554540433Z',
    message: _.uniqueId('message'),
  };
}

describe('LogsComponent <-> PodsLogsService', () => {
  let component: LogsComponent;
  let fixture: ComponentFixture<LogsComponent>;
  let service: PodLogsService;
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LogsComponent],
      providers: [PodLogsService],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    service = TestBed.get(PodLogsService);
    httpClient = TestBed.get(HttpClient);
    httpTestingController = TestBed.get(HttpTestingController);
  });

  it('should render list of containers to choose from', () => {
    expect(component.selectedContainer).toBe('');

    component.view = createTestLogsView([
      'containerA',
      'containerB',
      'containerC',
    ]);
    fixture.detectChanges();

    const selectOptionsDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.container-select select > option')
    );
    expect(selectOptionsDebugElements.length).toBe(3);
    expect(selectOptionsDebugElements[0].nativeElement.value).toBe(
      'containerA'
    );
    expect(selectOptionsDebugElements[1].nativeElement.value).toBe(
      'containerB'
    );
    expect(selectOptionsDebugElements[2].nativeElement.value).toBe(
      'containerC'
    );
  });

  it('should allow user to toggle displaying timestamps', () => {
    component.shouldDisplayTimestamp = true;
    component.containerLogs = [
      { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageA' },
      { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageB' },
      { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageC' },
    ];
    fixture.detectChanges();

    let logEntriesDebugElement: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.container-log')
    );
    expect(logEntriesDebugElement.length).toBe(3);
    expect(logEntriesDebugElement[0].nativeElement.textContent).toMatch(
      /May \d+, 2019(.+)messageA/
    );
    expect(logEntriesDebugElement[1].nativeElement.textContent).toMatch(
      /May \d+, 2019(.*)+messageB/
    );
    expect(logEntriesDebugElement[2].nativeElement.textContent).toMatch(
      /May \d+, 2019(.*)+messageC/
    );

    component.shouldDisplayTimestamp = false;
    fixture.detectChanges();

    logEntriesDebugElement = fixture.debugElement.queryAll(
      By.css('.container-log')
    );
    expect(logEntriesDebugElement.length).toBe(3);
    expect(logEntriesDebugElement[0].nativeElement.textContent).toMatch(
      /^\s+messageA\s+$/
    );
    expect(logEntriesDebugElement[1].nativeElement.textContent).toMatch(
      /^\s+messageB\s+$/
    );
    expect(logEntriesDebugElement[2].nativeElement.textContent).toMatch(
      /^\s+messageC\s+$/
    );
  });

  it('should continously scroll to new logs if user has already scrolled to the bottom', () => {
    const numberOfEntriesRequiredToScroll = 200;
    component.containerLogs = _.map(
      _.range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    fixture.detectChanges();

    const logWrapperDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    );
    let logWrapperNativeElement: HTMLDivElement =
      logWrapperDebugElement.nativeElement;
    expect(logWrapperNativeElement.scrollHeight).toBeGreaterThan(
      logWrapperNativeElement.clientHeight
    );
    expect(logWrapperNativeElement.scrollTop).toBe(0);

    const logWrapperHeight = logWrapperNativeElement.clientHeight;
    logWrapperNativeElement.scrollTop = logWrapperNativeElement.scrollHeight;

    expect(logWrapperNativeElement.scrollTop).toEqual(
      logWrapperNativeElement.scrollHeight - logWrapperHeight
    );

    const newContainerLogs: LogEntry[] = _.map(
      _.range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    component.containerLogs.push(...newContainerLogs);
    logWrapperNativeElement.dispatchEvent(new Event('scroll'));
    fixture.detectChanges();

    logWrapperNativeElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    ).nativeElement;
    expect(logWrapperNativeElement.scrollTop).toEqual(
      logWrapperNativeElement.scrollHeight - logWrapperHeight
    );
  });

  it('should keep scroll position even if new logs are coming in and user is not at bottom', () => {
    const numberOfEntriesRequiredToScroll = 200;
    component.containerLogs = _.map(
      _.range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    fixture.detectChanges();

    const logWrapperDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    );
    let logWrapperNativeElement: HTMLDivElement =
      logWrapperDebugElement.nativeElement;
    expect(logWrapperNativeElement.scrollHeight).toBeGreaterThan(
      logWrapperNativeElement.clientHeight
    );
    expect(logWrapperNativeElement.scrollTop).toBe(0);

    // scroll halfway
    const halfwayScrollMark = Math.floor(
      logWrapperNativeElement.clientHeight / 2
    );
    logWrapperNativeElement.scrollTop = halfwayScrollMark;
    logWrapperNativeElement.dispatchEvent(new Event('scroll'));

    // add new logs
    const newContainerLogs: LogEntry[] = _.map(
      _.range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    component.containerLogs.push(...newContainerLogs);
    fixture.detectChanges();

    // check scroll is in same place
    logWrapperNativeElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    ).nativeElement;
    expect(logWrapperNativeElement.scrollTop).toBe(halfwayScrollMark);
  });

  afterEach(() => {
    httpTestingController.verify();
  });
});
