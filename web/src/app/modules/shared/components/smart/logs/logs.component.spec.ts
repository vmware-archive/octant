// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { LogsComponent } from './logs.component';
import { LogEntry } from 'src/app/modules/shared/models/content';
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
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should stay at the bottom of the container when new logs arrive', () => {
    const { nativeElement } = component.scrollTarget;
    component.scrollToBottom = true;

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
    expect(selectHighlights.length).toEqual(15);
    expect(selectHighlights[0].nativeElement.innerText).toEqual('Just');
  });
});
