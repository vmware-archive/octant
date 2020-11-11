// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { HttpClient } from '@angular/common/http';
import { HttpTestingController } from '@angular/common/http/testing';
import { DebugElement } from '@angular/core';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import map from 'lodash/map';
import range from 'lodash/range';
import uniqueId from 'lodash/uniqueId';
import {
  LogEntry,
  LogsView,
  Since,
} from 'src/app/modules/shared/models/content';
import { PodLogsService } from 'src/app/modules/shared/pod-logs/pod-logs.service';
import { LogsComponent } from './logs.component';
import { AnsiPipe } from '../../../pipes/ansiPipe/ansi.pipe';
import { windowProvider, WindowToken } from '../../../../../window';
import { StringEscapePipe } from '../../../pipes/stringEscape/string.escape.pipe';

function createTestLogsView(
  durations: Since[],
  containers: string[]
): LogsView {
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
      durations,
    },
  };
}

function createRandomLogEntry(): LogEntry {
  return {
    timestamp: '2019-05-06T18:59:06.554540433Z',
    message: uniqueId('message'),
    container: 'test-container',
  };
}

describe('LogsComponent <-> PodsLogsService', () => {
  let component: LogsComponent;
  let fixture: ComponentFixture<LogsComponent>;
  let service: PodLogsService;
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  const defaultTestLogs = [
    {
      timestamp: '2019-05-06T18:59:06.554540433Z',
      message: 'messageA',
      container: 'test-container',
    },
    {
      timestamp: '2019-05-06T18:59:06.554540433Z',
      message: 'messageB',
      container: 'test-container',
    },
    {
      timestamp: '2019-05-06T18:59:06.554540433Z',
      message: 'messageC',
      container: 'test-container',
    },
  ];

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LogsComponent, AnsiPipe, StringEscapePipe],
      providers: [
        PodLogsService,
        { provide: WindowToken, useFactory: windowProvider },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LogsComponent);
    component = fixture.componentInstance;
    service = TestBed.inject(PodLogsService);
    httpClient = TestBed.inject(HttpClient);
    httpTestingController = TestBed.inject(HttpTestingController);

    component.view = createTestLogsView(
      [{ label: '5 minutes', seconds: 300 }],
      ['containerA', 'containerB', 'containerC']
    );
  });

  afterEach(() => {
    TestBed.resetTestingModule();
  });

  it('should allow user to toggle displaying timestamps', () => {
    component.shouldDisplayTimestamp = true;
    component.containerLogs = defaultTestLogs;
    fixture.detectChanges();

    let logEntriesDebugElement: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.container-log')
    );
    expect(logEntriesDebugElement.length).toBe(3);
    expect(logEntriesDebugElement[0].nativeElement.textContent).toMatch(
      /May \d+, 2019(.+)messageA/
    );
    expect(logEntriesDebugElement[1].nativeElement.textContent).toMatch(
      /May \d+, 2019(.+)messageB/
    );
    expect(logEntriesDebugElement[2].nativeElement.textContent).toMatch(
      /May \d+, 2019(.+)messageC/
    );

    component.shouldDisplayTimestamp = false;
    fixture.detectChanges();

    logEntriesDebugElement = fixture.debugElement.queryAll(
      By.css('.container-log')
    );
    expect(logEntriesDebugElement.length).toBe(3);
    expect(logEntriesDebugElement[0].nativeElement.textContent).toBe(
      'test-containermessageA'
    );
    expect(logEntriesDebugElement[1].nativeElement.textContent).toBe(
      'test-containermessageB'
    );
    expect(logEntriesDebugElement[2].nativeElement.textContent).toBe(
      'test-containermessageC'
    );
  });

  it('should continuously scroll to new logs if user has already scrolled to the bottom', () => {
    const numberOfEntriesRequiredToScroll = 200;
    component.containerLogs = map(
      range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    fixture.detectChanges();

    const logWrapperDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.log-container')
    );
    let logWrapperNativeElement: HTMLDivElement =
      logWrapperDebugElement.nativeElement;
    expect(logWrapperNativeElement.scrollHeight).toEqual(
      logWrapperNativeElement.clientHeight
    );
    expect(logWrapperNativeElement.scrollTop).toEqual(0);

    const logWrapperHeight = logWrapperNativeElement.clientHeight;
    logWrapperNativeElement.scrollTop = logWrapperNativeElement.scrollHeight;

    expect(logWrapperNativeElement.scrollTop).toEqual(
      logWrapperNativeElement.scrollHeight - logWrapperHeight
    );

    const newContainerLogs: LogEntry[] = map(
      range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    component.containerLogs.push(...newContainerLogs);
    logWrapperNativeElement.dispatchEvent(new Event('scroll'));
    fixture.detectChanges();

    logWrapperNativeElement = fixture.debugElement.query(
      By.css('.log-container')
    ).nativeElement;
    expect(logWrapperNativeElement.scrollTop).toEqual(0);
  });

  it('should keep scroll position even if new logs are coming in and user is not at bottom', () => {
    const numberOfEntriesRequiredToScroll = 200;
    component.containerLogs = map(
      range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    fixture.detectChanges();

    const logWrapperDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    );
    let logWrapperNativeElement: HTMLDivElement =
      logWrapperDebugElement.nativeElement;
    fixture.whenStable().then(() => {
      expect(logWrapperNativeElement.scrollHeight).toBeGreaterThan(
        logWrapperNativeElement.clientHeight
      );
      expect(logWrapperNativeElement.scrollTop).toBeGreaterThan(0);
    });

    // scroll halfway
    const halfwayScrollMark = Math.floor(
      logWrapperNativeElement.clientHeight / 2
    );
    logWrapperNativeElement.scrollTop = halfwayScrollMark;
    logWrapperNativeElement.dispatchEvent(new Event('scroll'));

    // add new logs
    const newContainerLogs: LogEntry[] = map(
      range(numberOfEntriesRequiredToScroll),
      createRandomLogEntry
    );
    component.containerLogs.push(...newContainerLogs);
    fixture.detectChanges();

    // check scroll is in same place
    logWrapperNativeElement = fixture.debugElement.query(
      By.css('.container-logs-bg')
    ).nativeElement;
    fixture.whenStable().then(() => {
      expect(logWrapperNativeElement.scrollTop).toBe(halfwayScrollMark);
    });
  });

  it('should filter messages based on search string', () => {
    component.view = createTestLogsView(
      [{ label: '5 minutes', seconds: 300 }],
      ['containerA', 'containerB', 'containerC']
    );

    component.shouldDisplayTimestamp = true;
    component.containerLogs = defaultTestLogs;
    component.filterText = 'message';
    fixture.detectChanges();
    VerifyElementsExist('.container-log', 3);
    VerifyElementsExist('.highlight', 3);
    VerifyElementsExist('.highlight-selected', 1);

    component.filterText = 'messageA';
    fixture.detectChanges();
    VerifyElementsExist('.container-log', 3);
    VerifyElementsExist('.highlight', 1);
    VerifyElementsExist('.highlight-selected', 1);

    component.showOnlyFiltered = true;
    fixture.detectChanges();
    VerifyElementsExist('.container-log', 1);
    VerifyElementsExist('.highlight', 1);
    VerifyElementsExist('.highlight-selected', 1);

    component.filterText = '';
    fixture.detectChanges();
    VerifyElementsExist('.container-log', 3);
    VerifyElementsExist('.highlight', 0);
    VerifyElementsExist('.highlight-selected', 0);
  });

  afterEach(() => {
    httpTestingController.verify();
  });

  function VerifyElementsExist(selector: string, noItems: number) {
    const element: DebugElement[] = fixture.debugElement.queryAll(
      By.css(selector)
    );
    expect(element.length).toBe(noItems);
  }
});
