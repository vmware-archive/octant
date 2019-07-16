// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  async,
  ComponentFixture,
  TestBed,
  fakeAsync,
  tick,
} from '@angular/core/testing';
import { PortsComponent } from './ports.component';
import { PortForwardService } from 'src/app/services/port-forward/port-forward.service';
import { HttpClient } from '@angular/common/http';
import { HttpTestingController } from '@angular/common/http/testing';
import { Port, PortsView } from 'src/app/models/content';
import _ from 'lodash';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import getAPIBase from 'src/app/services/common/getAPIBase';
import { notifierServiceStubFactory } from 'src/app/testing/notifier-service.stub';
import {
  NotifierService,
  NotifierSignalType,
} from 'src/app/services/notifier/notifier.service';

const API_BASE = getAPIBase();

function createTestPortsView(ports: Port[]): PortsView {
  return {
    metadata: { type: 'ports', title: [], accessor: '' },
    config: { ports },
  };
}

function createTestPort(
  state: Partial<{ port: number; isForwardable: boolean; isForwarded: boolean }>
): Port {
  return {
    metadata: { type: 'port', title: [], accessor: '' },
    config: {
      namespace: 'default',
      apiVersion: 'v1',
      kind: 'Pod',
      name: 'cartservice-pod',
      port: 8080,
      protocol: 'TCP',
      state: _.assign(
        {
          id: _.uniqueId('portforward-'),
        },
        state
      ),
    },
  };
}

describe('PortForwardComponent <-> PortForwardService', () => {
  let component: PortsComponent;
  let fixture: ComponentFixture<PortsComponent>;
  let service: PortForwardService;
  let httpClient: HttpClient;
  let httpTestingController: HttpTestingController;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [PortsComponent],
      providers: [
        { provide: NotifierService, useFactory: notifierServiceStubFactory },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PortsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
    service = TestBed.get(PortForwardService);
    httpClient = TestBed.get(HttpClient);
    httpTestingController = TestBed.get(HttpTestingController);
  });

  it('renders rows of portforwarding options', () => {
    const testPorts: Port[] = [
      createTestPort({ isForwardable: true }),
      createTestPort({ isForwarded: true, port: 12345 }),
      createTestPort({ isForwardable: false }),
      createTestPort({ isForwarded: false }),
    ];
    const testPortsView = createTestPortsView(testPorts);
    component.view = testPortsView;
    fixture.detectChanges();

    const forwardedLinkDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.open-pf')
    );
    expect(forwardedLinkDebugElements.length).toBe(1);
    expect(
      (forwardedLinkDebugElements[0].nativeElement as HTMLAnchorElement)
        .textContent
    ).toMatch(/localhost:12345/i);

    const forwardableDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.start-pf')
    );
    expect(forwardableDebugElements.length).toBe(1);
  });

  it('starts portforwarding on btn click', () => {
    const testPorts: Port[] = [createTestPort({ isForwardable: true })];
    const testPortsView = createTestPortsView(testPorts);
    component.view = testPortsView;
    fixture.detectChanges();

    const forwardableDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.start-pf')
    );
    expect(forwardableDebugElements.length).toBe(1);

    const startForwardButtonNativeElement: HTMLButtonElement =
      forwardableDebugElements[0].nativeElement;
    startForwardButtonNativeElement.dispatchEvent(new Event('click'));

    const req = httpTestingController.expectOne(
      `${API_BASE}/api/v1/content/overview/port-forwards`
    );
    expect(req.request.method).toBe('POST');

    req.flush({});
    httpTestingController.verify();
  });

  it('removes portforwarding on btn click', () => {
    const testPort = createTestPort({ isForwarded: true, port: 54321 });
    const testPorts: Port[] = [testPort];
    const testPortsView = createTestPortsView(testPorts);
    component.view = testPortsView;
    fixture.detectChanges();

    const forwardedLinkDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.remove-pf')
    );
    expect(forwardedLinkDebugElements.length).toBe(1);

    forwardedLinkDebugElements[0].nativeElement.dispatchEvent(
      new Event('click')
    );

    const req = httpTestingController.expectOne(
      `${API_BASE}/api/v1/content/overview/port-forwards/${
        testPort.config.state.id
      }`
    );
    expect(req.request.method).toBe('DELETE');

    req.flush({});
    httpTestingController.verify();
  });

  it('toggles between port forward button states', () => {
    const testPortsView = createTestPortsView([
      createTestPort({ isForwardable: true }),
      createTestPort({ isForwardable: true }),
      createTestPort({ isForwarded: true, port: 7777 }),
    ]);
    component.view = testPortsView;
    fixture.detectChanges();

    let portDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.port')
    );
    expect(portDebugElements.length).toBe(3);
    expect(
      portDebugElements[0].query(By.css('.port-text')).nativeElement.textContent
    ).toMatch(/8080\/TCP/);
    expect(portDebugElements[0].query(By.css('.start-pf'))).toBeTruthy();
    expect(portDebugElements[1].query(By.css('.start-pf'))).toBeTruthy();
    expect(portDebugElements[2].query(By.css('.remove-pf'))).toBeTruthy();

    _.assign(testPortsView.config.ports[0].config.state, {
      isForwardable: false,
      isForwarded: true,
      port: 44444,
    });
    _.assign(testPortsView.config.ports[1].config.state, {
      isForwardable: false,
      isForwarded: true,
      port: 56789,
    });
    _.assign(testPortsView.config.ports[2].config.state, {
      isForwardable: true,
      isForwarded: false,
    });
    fixture.detectChanges();

    portDebugElements = fixture.debugElement.queryAll(By.css('.port'));
    expect(portDebugElements.length).toBe(3);
    expect(
      portDebugElements[0].query(By.css('.open-pf')).nativeElement.textContent
    ).toMatch(/localhost:44444/i);
    expect(portDebugElements[0].query(By.css('.remove-pf'))).toBeTruthy();
    expect(
      portDebugElements[1].query(By.css('.open-pf')).nativeElement.textContent
    ).toMatch(/localhost:56789/i);
    expect(portDebugElements[2].query(By.css('.start-pf'))).toBeTruthy();
  });

  it('notifies users if there was an error starting the portforward', () => {
    const notifierService = TestBed.get(NotifierService);
    const { notifierSessionStub } = notifierService;

    const testPorts: Port[] = [createTestPort({ isForwardable: true })];
    const testPortsView = createTestPortsView(testPorts);
    component.view = testPortsView;
    fixture.detectChanges();

    const forwardableDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.start-pf')
    );
    expect(forwardableDebugElements.length).toBe(1);

    const startForwardButtonNativeElement: HTMLButtonElement =
      forwardableDebugElements[0].nativeElement;
    startForwardButtonNativeElement.dispatchEvent(new Event('click'));

    const req = httpTestingController.expectOne(
      `${API_BASE}/api/v1/content/overview/port-forwards`
    );
    expect(req.request.method).toBe('POST');

    req.flush({}, { status: 500, statusText: 'Error' });

    expect((notifierSessionStub.pushSignal as jasmine.Spy).calls.count()).toBe(
      1
    );
    expect(
      (notifierSessionStub.pushSignal as jasmine.Spy).calls.first().args[0]
    ).toBe(NotifierSignalType.ERROR);

    httpTestingController.verify();
  });
});
