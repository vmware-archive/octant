/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';

import {
  ContentPathUpdate,
  ContentPathUpdateMessage,
  ContentService,
  ContentUpdateMessage,
} from './content.service';
import { WebsocketServiceMock } from '../websocket/mock';
import {
  BackendService,
  WebsocketService,
} from '../websocket/websocket.service';
import { Content } from '../../../../models/content';
import { Router } from '@angular/router';
import {
  Filter,
  LabelFilterService,
} from '../../../../services/label-filter/label-filter.service';

describe('ContentService', () => {
  let service: ContentService;
  const mockRouter = {
    navigate: jasmine.createSpy('navigate'),
  };

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
      ],
    });

    service = TestBed.get(ContentService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('content update', () => {
    const update: Content = { title: [], viewComponents: [] };

    beforeEach(() => {
      const backendService = TestBed.get(WebsocketService);
      backendService.triggerHandler(ContentUpdateMessage, update);
    });

    it('triggers a content change', () => {
      service.current.subscribe(current =>
        expect(current).toEqual({ content: update })
      );
    });
  });

  describe('content path update', () => {
    let backendService: BackendService;

    beforeEach(() => {
      backendService = TestBed.get(WebsocketService);
    });

    describe('without query params', () => {
      const update: ContentPathUpdate = {
        contentPath: 'path',
        queryParams: {},
      };

      beforeEach(() => {
        backendService.triggerHandler(ContentPathUpdateMessage, update);
      });

      it('triggers a content change', () => {
        expect(mockRouter.navigate).toHaveBeenCalledWith(['content', 'path'], {
          queryParams: {},
        });
      });
    });

    describe('with query params', () => {
      const update: ContentPathUpdate = {
        contentPath: 'path',
        queryParams: {
          foo: ['bar'],
        },
      };

      beforeEach(() => {
        backendService.triggerHandler(ContentPathUpdateMessage, update);
      });

      it('triggers a content change', () => {
        expect(mockRouter.navigate).toHaveBeenCalledWith(['content', 'path'], {
          queryParams: {
            foo: ['bar'],
          },
        });
      });
    });
  });

  describe('label filters updated', () => {
    let labelFilterService: LabelFilterService;

    const filters = [{ key: 'foo', value: 'bar' }];

    beforeEach(() => {
      labelFilterService = TestBed.get(LabelFilterService);
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
      backendService = TestBed.get(WebsocketService);
      spyOn(backendService, 'sendMessage');
    });

    it('sends a setContentPath message to the server', () => {
      service.setContentPath('path');
      expect(backendService.sendMessage).toHaveBeenCalledWith(
        'setContentPath',
        {
          contentPath: 'path',
          filters: [],
        }
      );
    });

    describe('with filters defined', () => {
      beforeEach(() => {
        filters = [{ key: 'foo', value: 'bar' }];
        const labelFilterService = TestBed.get(LabelFilterService);
        labelFilterService.filters.next(filters);
      });

      it('sends a setContentPath message to the server', () => {
        service.setContentPath('path');
        expect(backendService.sendMessage).toHaveBeenCalledWith(
          'setContentPath',
          {
            contentPath: 'path',
            filters,
          }
        );
      });
    });
  });

  describe('set query params', () => {
    let backendService: BackendService;

    beforeEach(() => {
      backendService = TestBed.get(WebsocketService);
      spyOn(backendService, 'sendMessage');
    });

    it('sends a setQueryParams message to the server', () => {
      service.setQueryParams({ foo: 'bar' });
      expect(backendService.sendMessage).toHaveBeenCalledWith(
        'setQueryParams',
        {
          params: { foo: 'bar' },
        }
      );
    });
  });
});
