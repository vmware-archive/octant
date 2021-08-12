/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';

import {
  ContentService,
  ContentUpdate,
  ContentUpdateMessage,
} from './content.service';
import { WebsocketServiceMock } from '../../../../data/services/websocket/mock';
import {
  BackendService,
  WebsocketService,
} from '../../../../data/services/websocket/websocket.service';
import { Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../label-filter/label-filter.service';
import { Title } from '@angular/platform-browser';

describe('ContentService', () => {
  let service: ContentService;
  const mockRouter = jasmine.createSpyObj('Router', ['navigate']);

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ContentService,
        LabelFilterService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
        {
          provide: Router,
          useValue: mockRouter,
        },
        {
          provide: Title,
          useValue: jasmine.createSpyObj('Title', ['getTitle', 'setTitle']),
        },
      ],
    });

    service = TestBed.inject(ContentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('content update', () => {
    const update: ContentUpdate = {
      content: { extensionComponent: null, title: [], viewComponents: [] },
      namespace: 'default',
      contentPath: '/path',
      queryParams: {},
    };

    beforeEach(() => {
      const backendService = TestBed.inject(WebsocketService);
      backendService.triggerHandler(ContentUpdateMessage, update);
    });

    it('triggers a content change', () => {
      service.current.subscribe(current =>
        expect(current).toEqual({
          content: update.content,
          currentPath: '/path',
        })
      );
    });
  });

  describe('label filters updated', () => {
    let labelFilterService: LabelFilterService;

    const filters = [{ key: 'foo', value: 'bar' }];

    beforeEach(() => {
      labelFilterService = TestBed.inject(LabelFilterService);
      labelFilterService.filters.next(filters);
    });

    it('updates local filters', () => {
      expect(service.currentFilters).toEqual(filters);
    });
  });

  describe('set content path', () => {
    let backendService: BackendService;
    let filters: Filter[];

    beforeEach(() => {
      backendService = TestBed.inject(WebsocketService);
      spyOn(backendService, 'sendMessage');
    });

    it('sends a setContentPath message to the server', () => {
      service.setContentPath('path', {});
      expect(backendService.sendMessage).toHaveBeenCalledWith(
        'action.octant.dev/setContentPath',
        {
          contentPath: 'path',
          params: {},
        }
      );
    });

    describe('with filters defined', () => {
      beforeEach(() => {
        filters = [{ key: 'foo', value: 'bar' }];
        const labelFilterService = TestBed.inject(LabelFilterService);
        labelFilterService.filters.next(filters);
      });

      it('sends a setContentPath message to the server', () => {
        service.setContentPath('path', { filters });
        expect(backendService.sendMessage).toHaveBeenCalledWith(
          'action.octant.dev/setContentPath',
          {
            contentPath: 'path',
            params: { filters },
          }
        );
      });
    });
  });
});
