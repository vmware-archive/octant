// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { LogsComponent } from './logs.component';
import { LogEntry, LogsView } from 'src/app/modules/shared/models/content';
import { By } from '@angular/platform-browser';
import { DebugElement, ElementRef } from '@angular/core';
import { AnsiPipe } from '../../../pipes/ansiPipe/ansi.pipe';
import { windowProvider, WindowToken } from '../../../../../window';

/**
 * Adds lines of logs to LogsComponent
 * @param currentLogList Current list of logs when this function is called.
 * @param lines Number of log lines to add
 * @returns New list of logs.
 */
const addLogsToList = (
  currentLogList: LogEntry[],
  lines: number
): LogEntry[] => {
  const logList = [...currentLogList];
  for (let i = 1; i <= lines; i++) {
    logList.push({
      timestamp: '2019-08-19T12:07:00.1222053Z',
      message: 'Just for test',
      container: 'test-container',
    });
  }

  return logList;
};

describe('LogsComponent', () => {
  let component: LogsComponent;
  let fixture: ComponentFixture<LogsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [LogsComponent, AnsiPipe],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(LogsComponent);
    component = fixture.componentInstance;
    component.view = {
      metadata: {
        type: 'logs',
        title: [],
        accessor: 'logs',
      },
      config: {
        namespace: 'default',
        name: 'cartpod',
        containers: [],
        durations: [],
      },
    } as LogsView;

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should stay at the bottom of the container when new logs arrive', () => {
    const scrolltarget: ElementRef = component.scrollTarget;
    component.containerLogs = addLogsToList([], 15);
    fixture.detectChanges();
    expect(component.containerLogs.length).toBe(15);

    // This is done via CSS but scrollTop is unreliable so checking the class is enabled
    expect(scrolltarget.nativeElement.classList).toContain('container-content');
  });

  it('should filter messages based on regex expression', () => {
    component.filterText = '([A-Z])\\w+';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([], 30);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(150);
    expect(selectHighlights[0].nativeElement.innerText).toEqual('test');
  });

  it('should filter for positive lookahead regex', () => {
    component.filterText = '(?=Just)';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([], 30);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(30);
    expect(selectHighlights[0].nativeElement.innerText).toEqual(
      'Just for test'
    );
  });

  it('should filter case insensitive', () => {
    component.filterText = 'JUST';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([], 30);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(30);
    expect(selectHighlights[0].nativeElement.innerText).toEqual('Just');
  });

  it('forward button should wrap search at bottom', () => {
    component.containerLogs = [
      {
        timestamp: '2019-05-06T18:50:06.554540433Z',
        message: 'Test log line 1',
        container: 'test-container',
      },
      {
        timestamp: '2019-05-06T18:59:06.554540433Z',
        message: 'Test log line 2',
        container: 'test-container',
      },
    ];
    component.filterText = 'Test log';
    fixture.detectChanges();

    const prevButton = fixture.debugElement.nativeElement.querySelector(
      '#button-prev'
    );
    const nextButton = fixture.debugElement.nativeElement.querySelector(
      '#button-next'
    );
    const badgeElement: HTMLDivElement = fixture.debugElement.query(
      By.css('.clr-filter-summary')
    ).nativeElement;

    expect(badgeElement.innerText).toBe('1/2 items');
    nextButton.click();

    fixture.whenStable().then(() => {
      const offsetSecondElement = getSelectedHighlightTop();

      fixture.detectChanges();
      expect(badgeElement.innerText).toBe('2/2 items');

      nextButton.click();
      fixture.detectChanges();
      expect(getSelectedHighlightTop()).toBeLessThan(offsetSecondElement); // should roll-up to 1st
      expect(badgeElement.innerText).toBe('1/2 items');

      prevButton.click();
      fixture.detectChanges();
      expect(getSelectedHighlightTop()).toBe(offsetSecondElement); // should come back to 2nd
      expect(badgeElement.innerText).toBe('2/2 items');
    });
  });

  function getSelectedHighlightTop() {
    const nextSelectedElement: HTMLDivElement = fixture.debugElement.query(
      By.css('.highlight-selected')
    ).nativeElement;

    return nextSelectedElement.offsetTop;
  }
});
