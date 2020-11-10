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
import { StringEscapePipe } from '../../../pipes/stringEscape/string.escape.pipe';

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
        declarations: [LogsComponent, AnsiPipe, StringEscapePipe],
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

  it('should handle text that has to be escaped (#1582)', () => {
    const expectedLogEntry = `2020-11-01T11:56:45.910889701Z {"@timestamp":"2020-11-01T11:56:45.909Z","level":"INFO","class":"rs.jms.message.received.from.mq","message":"GenericMessage [payload=<?xml version="1.0" encoding="UTF-8"?><ns:Message\\n xmlns:ns="http://otpbank.ru/Message"\\n><MessageHeader\\n><MessageDate\\n>2020-11-01T14:56:45</MessageDate\\n><ProcessID\\n>Broker</ProcessID\\n><MessageType\\n>CreateOpty</MessageType\\n><MessageID\\n>100B2F2E5E64EDB65E6E0533CA0A8C03EDD</MessageID\\n><IntegrationID\\n>4c020f7acfd84c2d8af05dfbd1404e1c</IntegrationID\\n><InitiatorSystem\\n>SIEBEL</InitiatorSystem\\n><SourceSystem\\n>SIEBEL</SourceSystem\\n><TargetSystemList\\n><TargetSystem\\n>mcs-posbroker-api</TargetSystem\\n></TargetSystemList\\n><Version\\n>1.0</Version\\n><ResultInfo\\n><Code\\n>0</Code\\n><Description\\n></Description\\n><Exception\\n></Exception\\n></ResultInfo\\n></MessageHeader\\n><MessageBody\\n><CreateOptyResponse\\n><Opty_Id\\n>2-1OGXGXFW</Opty_Id\\n><Error_Code\\n>0</Error_Code\\n><Error_Message\\n></Error_Message\\n></CreateOptyResponse\\n></MessageBody\\n></ns:Message\\n>, headers={JMS_IBM_Character_Set=UTF-8, errorChannel=org.springframework.messaging.core.GenericMessagingTemplate$TemporaryReplyChannel@b8364f77, JMS_IBM_MsgType=8, jms_destination=queue://WMB02PQM/POS_CBCREATEOPTY_SIEBEL_RES, JMSXUserID=siebel , JMS_IBM_Encoding=273, priority=4, jms_timestamp=1604231805837, JMSXAppID=WebSphere MQ Client for Java, JMS_IBM_PutApplType=28, JMS_IBM_Format=MQSTR , replyChannel=org.springframework.messaging.core.GenericMessagingTemplate$TemporaryReplyChannel@b8364f77, jms_redelivered=false, JMS_IBM_PutDate=20201101, JMSXDeliveryCount=1, jms_correlationId=784b0c72-8937-4615-b911-54b656e82460, ws_soapAction="document/http://siebel.com/CustomUI:CreateOpty_spcv2", JMS_IBM_PutTime=11564589, jms_type=SiebelJMSMessage, id=fa092de7-e6c7-a00e-3647-d3c6a8eae2fe, jms_messageId=ID:414d5120574d42303250514d20202020667e9e5f298f3623, timestamp=1604231805909}]","trace":"","span":"","parent":"","thread":"http-nio-8080-exec-9"}`;
    const logList = [];
    logList.push({
      timestamp: '2020-11-01T11:56:45.910889701Z',
      message: `2020-11-01T11:56:45.910889701Z {"@timestamp":"2020-11-01T11:56:45.909Z","level":"INFO","class":"rs.jms.message.received.from.mq","message":"GenericMessage [payload=<?xml version=\"1.0\" encoding=\"UTF-8\"?><ns:Message\n xmlns:ns=\"http://otpbank.ru/Message\"\n><MessageHeader\n><MessageDate\n>2020-11-01T14:56:45</MessageDate\n><ProcessID\n>Broker</ProcessID\n><MessageType\n>CreateOpty</MessageType\n><MessageID\n>100B2F2E5E64EDB65E6E0533CA0A8C03EDD</MessageID\n><IntegrationID\n>4c020f7acfd84c2d8af05dfbd1404e1c</IntegrationID\n><InitiatorSystem\n>SIEBEL</InitiatorSystem\n><SourceSystem\n>SIEBEL</SourceSystem\n><TargetSystemList\n><TargetSystem\n>mcs-posbroker-api</TargetSystem\n></TargetSystemList\n><Version\n>1.0</Version\n><ResultInfo\n><Code\n>0</Code\n><Description\n></Description\n><Exception\n></Exception\n></ResultInfo\n></MessageHeader\n><MessageBody\n><CreateOptyResponse\n><Opty_Id\n>2-1OGXGXFW</Opty_Id\n><Error_Code\n>0</Error_Code\n><Error_Message\n></Error_Message\n></CreateOptyResponse\n></MessageBody\n></ns:Message\n>, headers={JMS_IBM_Character_Set=UTF-8, errorChannel=org.springframework.messaging.core.GenericMessagingTemplate$TemporaryReplyChannel@b8364f77, JMS_IBM_MsgType=8, jms_destination=queue://WMB02PQM/POS_CBCREATEOPTY_SIEBEL_RES, JMSXUserID=siebel , JMS_IBM_Encoding=273, priority=4, jms_timestamp=1604231805837, JMSXAppID=WebSphere MQ Client for Java, JMS_IBM_PutApplType=28, JMS_IBM_Format=MQSTR , replyChannel=org.springframework.messaging.core.GenericMessagingTemplate$TemporaryReplyChannel@b8364f77, jms_redelivered=false, JMS_IBM_PutDate=20201101, JMSXDeliveryCount=1, jms_correlationId=784b0c72-8937-4615-b911-54b656e82460, ws_soapAction=\"document/http://siebel.com/CustomUI:CreateOpty_spcv2\", JMS_IBM_PutTime=11564589, jms_type=SiebelJMSMessage, id=fa092de7-e6c7-a00e-3647-d3c6a8eae2fe, jms_messageId=ID:414d5120574d42303250514d20202020667e9e5f298f3623, timestamp=1604231805909}]","trace":"","span":"","parent":"","thread":"http-nio-8080-exec-9"}`,
      container: 'test-container',
    });

    component.shouldDisplayTimestamp = false;
    component.containerLogs = logList;
    fixture.detectChanges();

    const logLines: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.container-log-message')
    );
    expect(logLines.length).toEqual(1);
    expect(logLines[0].nativeElement.innerText).toEqual(expectedLogEntry);
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
