// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { LogsComponent } from './logs.component';
import { LogEntry, LogsView } from 'src/app/modules/shared/models/content';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { AnsiPipe } from '../../../pipes/ansiPipe/ansi.pipe';

/**
 * Adds 15 logs to the provided list.
 * @param currentLogList Current list of logs when this function is called.
 * @returns New list of logs.
 */
const addLogsToList = (currentLogList: LogEntry[]): LogEntry[] => {
  const logList = [...currentLogList];
  for (let i = 1; i <= 15; i++) {
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

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [LogsComponent, AnsiPipe],
    }).compileComponents();
  }));

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
    const { nativeElement } = component.scrollTarget;

    component.containerLogs = addLogsToList([]);
    fixture.detectChanges();
    expect(component.containerLogs.length).toBe(15);
    expect(nativeElement.scrollTop).toEqual(
      nativeElement.scrollHeight - nativeElement.offsetHeight
    );

    component.containerLogs = addLogsToList(component.containerLogs);
    fixture.detectChanges();
    expect(component.containerLogs.length).toBe(30);
    expect(nativeElement.scrollTop).toEqual(
      nativeElement.scrollHeight - nativeElement.offsetHeight
    );
  });

  it('should filter messages based on regex expression', () => {
    component.filterText = '([A-Z])\\w+';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([]);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(75);
    expect(selectHighlights[0].nativeElement.innerText).toEqual('test');
  });

  it('should filter for positive lookahead regex', () => {
    component.filterText = '(?=Just)';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([]);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(15);
    expect(selectHighlights[0].nativeElement.innerText).toEqual(
      'Just for test'
    );
  });

  it('should filter case insensitive', () => {
    component.filterText = 'JUST';
    component.shouldDisplayTimestamp = false;
    component.containerLogs = addLogsToList([]);
    fixture.detectChanges();

    const selectHighlights: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.highlight')
    );
    expect(selectHighlights.length).toEqual(15);
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
