// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed, fakeAsync, tick } from '@angular/core/testing';
import { NamespaceService } from './namespace.service';
import {
  NotifierService,
  NotifierSignalType,
} from '../notifier/notifier.service';
import {
  Streamer,
  ContentStreamService,
} from '../content-stream/content-stream.service';
import _ from 'lodash';
import { Router } from '@angular/router';
import { RouterTestingModule } from '@angular/router/testing';
import { NgZone } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Navigation } from 'src/app/models/navigation';

describe('NamespaceService', () => {
  let service: NamespaceService;
  let router: Router;
  let ngZone: NgZone;

  beforeEach(() => {
    const emptyNavigation: Navigation = {
      sections: [],
    };

    const contentStreamServiceStub = {
      namespaces: new BehaviorSubject<string[]>([]),
      registerStreamer: (name: string, handler: Streamer) => {},
      streamer: () => new BehaviorSubject(emptyNavigation),
      navigation: new BehaviorSubject<Navigation>(emptyNavigation)
    };

    const notifierServiceStub = {
      notifierSessionStub: jasmine.createSpyObj([
        'removeAllSignals',
        'pushSignal',
      ]),
      createSession() {
        return this.notifierSessionStub;
      },
    };

    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule.withRoutes([
          { path: 'content/overview/namespace/:namespaceID', children: [] },
        ]),
      ],
      providers: [
        NamespaceService,
        { provide: ContentStreamService, useValue: contentStreamServiceStub },
        { provide: NotifierService, useValue: notifierServiceStub },
      ],
    });

    service = TestBed.get(NamespaceService);
    router = TestBed.get(Router);
    ngZone = TestBed.get(NgZone);
  });

  it('should be created', () => {
    service = TestBed.get(NamespaceService);
    expect(service).toBeDefined();
  });

  it('should set namespace and change route', fakeAsync(() => {
    ngZone.run(() => {
      service = TestBed.get(NamespaceService);
      router = TestBed.get(Router);
      service.setNamespace('default');
      tick();
      expect(service.current.getValue()).toBe('default');
      expect(router.url).toBe('/content/overview/namespace/default');
    });
  }));

  it('should set namespace to not default and still change route', fakeAsync(() => {
    ngZone.run(() => {
      service = TestBed.get(NamespaceService);
      router = TestBed.get(Router);

      service.setNamespace('default');
      tick();
      expect(service.current.getValue()).toBe('default');
      expect(router.url).toBe('/content/overview/namespace/default');

      service.setNamespace('kube-public');
      tick();
      expect(service.current.getValue()).toBe('kube-public');
      expect(router.url).toBe('/content/overview/namespace/kube-public');
    });
  }));

  it('should change namespace if incoming namespace is valid', fakeAsync(() => {
    ngZone.run(() => {
      service = TestBed.get(NamespaceService);
      router = TestBed.get(Router);
      service.list.next(['namespaceA', 'namespaceB', 'namespaceC']);
      router.navigate(['/content', 'overview', 'namespace', 'namespaceB']);
      tick();
      expect(service.current.getValue()).toBe('namespaceB');
    });
  }));

  it('should send signals if incoming namespace is invalid', fakeAsync(() => {
    ngZone.run(() => {
      service = TestBed.get(NamespaceService);
      router = TestBed.get(Router);
      service.list.next(['namespaceA', 'namespaceB', 'namespaceC']);
      const notifierSession = TestBed.get(NotifierService).notifierSessionStub;
      router.navigate(['/content', 'overview', 'namespace', 'testns']);
      tick();
      expect(notifierSession.removeAllSignals.calls.count()).toBe(1);
      expect(notifierSession.pushSignal.calls.count()).toBe(1);
      expect(notifierSession.pushSignal.calls.first().args[0]).toEqual(
        NotifierSignalType.ERROR
      );
    });
  }));
});
