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
import { LogEntry, LogsView } from 'src/app/modules/shared/models/content';
import getAPIBase from 'src/app/modules/shared/services/common/getAPIBase';
import { PodLogsService } from 'src/app/modules/shared/pod-logs/pod-logs.service';
import { LogsComponent } from './logs.component';
import { AnsiPipe } from '../../../pipes/ansiPipe/ansi.pipe';

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
    message: uniqueId('message'),
  };
}

describe('LogsComponent <-> PodsLogsService', () => {
  let component: LogsComponent;
  let fixture: ComponentFixture<LogsComponent>;
  let service: PodLogsService;
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;
  const defaultTestLogs = [
    { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageA' },
    { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageB' },
    { timestamp: '2019-05-06T18:59:06.554540433Z', message: 'messageC' },
  ];

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LogsComponent, AnsiPipe],
      providers: [PodLogsService],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    service = TestBed.inject(PodLogsService);
    httpClient = TestBed.inject(HttpClient);
    httpTestingController = TestBed.inject(HttpTestingController);
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
      'messageA'
    );
    expect(logEntriesDebugElement[1].nativeElement.textContent).toBe(
      'messageB'
    );
    expect(logEntriesDebugElement[2].nativeElement.textContent).toBe(
      'messageC'
    );
  });

  it('should continously scroll to new logs if user has already scrolled to the bottom', () => {
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
    expect(logWrapperNativeElement.scrollHeight).toBeGreaterThan(
      logWrapperNativeElement.clientHeight
    );
    expect(logWrapperNativeElement.scrollTop).toBeGreaterThan(0);

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
      By.css('.container-logs-bg')
    ).nativeElement;
    expect(logWrapperNativeElement.scrollTop).toEqual(
      logWrapperNativeElement.scrollHeight - logWrapperHeight
    );
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
    expect(logWrapperNativeElement.scrollHeight).toBeGreaterThan(
      logWrapperNativeElement.clientHeight
    );
    expect(logWrapperNativeElement.scrollTop).toBeGreaterThan(0);

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
    expect(logWrapperNativeElement.scrollTop).toBe(halfwayScrollMark);
  });

  it('should filer messages based on search string', () => {
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
