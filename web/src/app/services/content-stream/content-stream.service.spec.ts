// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed, ComponentFixture, inject } from '@angular/core/testing';
import { ContentStreamService } from './content-stream.service';
import { BehaviorSubject } from 'rxjs';
import { EventSourceStub, EventSourceService } from './event-source.service';
import {
  LabelFilterService,
  Filter,
} from '../label-filter/label-filter.service';
import {
  NotifierService,
  NotifierSignalType,
} from '../notifier/notifier.service';
import getAPIBase from '../common/getAPIBase';
import { ContentResponse } from '../../models/content';
import { Navigation } from '../../models/navigation';
import { notifierServiceStubFactory } from '../../testing/notifier-service.stub';
import { OverviewComponent } from '../../modules/overview/overview.component';
import { OverviewModule } from '../../modules/overview/overview.module';
import { NamespaceService } from '../namespace/namespace.service';
import { NamespaceComponent } from 'src/app/components/namespace/namespace.component';
import {
  KubeContextResponse,
  KubeContextService,
} from '../../modules/overview/services/kube-context/kube-context.service';
import { AppModule } from 'src/app/app.module';
import { NavigationComponent } from 'src/app/components/navigation/navigation.component';

const emptyContentResponse: ContentResponse = {
  content: {
    viewComponents: [],
    title: [],
  },
};

const emptyNavigation: Navigation = {
  sections: [],
};

const emptyKubeContext: KubeContextResponse = {
  contexts: [],
  currentContext: '',
};

describe('ContentStreamService', () => {
  const API_BASE = getAPIBase();
  let contentStreamService: ContentStreamService;
  let contextService: KubeContextService;
  let eventSourceService: {
    eventSourceStubs: Array<{ url: string; eventSourceStub: EventSourceStub }>;
  };
  let labelFilterService;
  let notifierService;
  let overviewFixture: ComponentFixture<OverviewComponent>;
  let namespaceFixture: ComponentFixture<NamespaceComponent>;
  let navigationFixture: ComponentFixture<NavigationComponent>;

  beforeEach(() => {
    const labelFilterStub: Partial<LabelFilterService> = {
      filters: new BehaviorSubject<Filter[]>([]),
    };

    const eventSourceServiceStub = {
      eventSourceStubs: [],
      createEventSource(url: string) {
        const eventSourceStub = new EventSourceStub();
        this.eventSourceStubs.push({ url, eventSourceStub });
        return eventSourceStub;
      },
    };

    TestBed.configureTestingModule({
      imports: [OverviewModule, AppModule],
      providers: [
        { provide: LabelFilterService, useValue: labelFilterStub },
        { provide: NotifierService, useFactory: notifierServiceStubFactory },
        { provide: EventSourceService, useValue: eventSourceServiceStub },
        NamespaceService,
        KubeContextService,
      ],
    }).compileComponents();

    contentStreamService = TestBed.get(ContentStreamService);
    contextService = TestBed.get(KubeContextService);
    eventSourceService = TestBed.get(EventSourceService);
    labelFilterService = TestBed.get(LabelFilterService);
    notifierService = TestBed.get(NotifierService);

    overviewFixture = TestBed.createComponent(OverviewComponent);
    overviewFixture.detectChanges();
    namespaceFixture = TestBed.createComponent(NamespaceComponent);
    namespaceFixture.detectChanges();
    navigationFixture = TestBed.createComponent(NavigationComponent);
    navigationFixture.detectChanges();
  });

  it('should create', () => {
    expect(contentStreamService).toBeTruthy();
  });

  it('should stream content after setting valid path w/o filters', () => {
    const { eventSourceStubs } = eventSourceService;
    const { notifierSessionStub } = notifierService;

    contentStreamService.openStream('namespace/default/overview');

    expect(notifierSessionStub.pushSignal.calls.count()).toBe(1);
    expect(notifierSessionStub.pushSignal.calls.first().args[0]).toBe(
      NotifierSignalType.LOADING
    );
    expect(eventSourceStubs.length).toBe(1);
    expect(eventSourceStubs[0].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5`
    );

    const { eventSourceStub } = eventSourceStubs[0];

    eventSourceStub.queueMessage(
      'content',
      JSON.stringify(emptyContentResponse)
    );
    eventSourceStub.queueMessage('navigation', JSON.stringify(emptyNavigation));
    eventSourceStub.queueMessage(
      'namespaces',
      JSON.stringify({ namespaces: [] })
    );
    eventSourceStub.flush();

    expect(contentStreamService.streamer('content').getValue()).toEqual(
      emptyContentResponse
    );
    expect(contentStreamService.streamer('navigation').getValue()).toEqual(
      emptyNavigation
    );
    expect(contentStreamService.streamer('namespaces').getValue()).toEqual([]);
    expect(contentStreamService.streamer('kubeConfig').getValue()).toEqual(
      emptyKubeContext
    );

    const testContentResponse: ContentResponse = {
      content: {
        title: [
          { metadata: { type: 'text', title: [], accessor: 'testTitle' } },
        ],
        viewComponents: [],
      },
    };

    eventSourceStub.queueMessage(
      'content',
      JSON.stringify(testContentResponse)
    );
    eventSourceStub.queueMessage('navigation', JSON.stringify(emptyNavigation));
    eventSourceStub.queueMessage(
      'namespaces',
      JSON.stringify({ namespaces: ['namespaceA', 'namespaceB'] })
    );
    eventSourceStub.flush();

    expect(contentStreamService.streamer('content').getValue()).toEqual(
      testContentResponse
    );
    expect(contentStreamService.streamer('navigation').getValue()).toEqual(
      emptyNavigation
    );
    expect(contentStreamService.streamer('namespaces').getValue()).toEqual([
      'namespaceA',
      'namespaceB',
    ]);
  });

  it('should stream content after setting valid path w/ filters', () => {
    const { eventSourceStubs } = eventSourceService;

    contentStreamService.openStream('namespace/default/overview');
    expect(eventSourceStubs.length).toBe(1);
    expect(eventSourceStubs[0].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5`
    );

    labelFilterService.filters.next([{ key: 'test1', value: 'value1' }]);

    expect(eventSourceStubs.length).toBe(2);
    expect(eventSourceStubs[1].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5&filter=test1%3Avalue1`
    );
  });

  it('should notify error signal if error is streamed in', () => {
    const { eventSourceStubs } = eventSourceService;
    const { notifierSessionStub } = notifierService;

    contentStreamService.openStream('namespace/default/overview');

    expect(eventSourceStubs.length).toBe(1);

    const { eventSourceStub } = eventSourceStubs[0];
    eventSourceStub.queueMessage('error');
    eventSourceStub.flush();

    expect(notifierSessionStub.pushSignal.calls.count()).toBe(2);
    expect(notifierSessionStub.pushSignal.calls.argsFor(1)[0]).toBe(
      NotifierSignalType.ERROR
    );
  });

  it('should notify warning signal if objectNotFound is streamed in', () => {
    const { eventSourceStubs } = eventSourceService;
    const { notifierSessionStub } = notifierService;

    contentStreamService.openStream('namespace/default/overview');

    expect(eventSourceStubs.length).toBe(1);

    const { eventSourceStub } = eventSourceStubs[0];
    eventSourceStub.queueMessage('objectNotFound', 'redirectpath');
    eventSourceStub.flush();

    expect(notifierSessionStub.pushSignal.calls.count()).toBe(3);
    expect(notifierSessionStub.pushSignal.calls.argsFor(2)[0]).toBe(
      NotifierSignalType.WARNING
    );
  });

  it('should cancel previous stream when setting up a new one', () => {
    const { eventSourceStubs } = eventSourceService;
    contentStreamService.openStream('namespace/default/overview');

    expect(eventSourceStubs.length).toBe(1);

    contentStreamService.openStream('namespace/testns/overview');

    expect(eventSourceStubs.length).toBe(2);
    expect(eventSourceStubs[1].url).toBe(
      `${API_BASE}/api/v1/content/namespace/testns/overview/?poll=5`
    );
  });

  it('should reset stream if filters change', () => {
    const { eventSourceStubs } = eventSourceService;
    contentStreamService.openStream('namespace/default/overview');

    expect(eventSourceStubs.length).toBe(1);
    expect(eventSourceStubs[0].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5`
    );

    labelFilterService.filters.next([{ key: 'test1', value: 'value1' }]);

    expect(eventSourceStubs.length).toBe(2);
    expect(eventSourceStubs[1].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5&filter=test1%3Avalue1`
    );

    labelFilterService.filters.next([
      { key: 'test1', value: 'value1' },
      { key: 'test2', value: 'value2' },
    ]);

    expect(eventSourceStubs.length).toBe(3);
    expect(eventSourceStubs[2].url).toBe(
      `${API_BASE}/api/v1/content/namespace/default/overview/?poll=5&filter=test1%3Avalue1&filter=test2%3Avalue2`
    );
  });
});
