// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { inject, TestBed } from '@angular/core/testing';
import { LabelFilterService } from './label-filter.service';
import { WebsocketService } from '../../../modules/overview/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../modules/overview/services/websocket/mock';

describe('LabelFilterService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        LabelFilterService,
        {
          provide: WebsocketService,
          useClass: WebsocketServiceMock,
        },
      ],
    });
  });

  describe('add', () => {
    it('should tell backend service to remove filter', inject(
      [LabelFilterService, WebsocketService],
      (service, websocketService) => {
        spyOn(websocketService, 'sendMessage');

        const filter = { key: 'foo', value: 'bar' };
        service.add(filter);

        expect(websocketService.sendMessage).toHaveBeenCalledWith('addFilter', {
          filter,
        });
      }
    ));
  });

  describe('remove', () => {
    it('should tell backend service to move filter', inject(
      [LabelFilterService, WebsocketService],
      (service, websocketService) => {
        spyOn(websocketService, 'sendMessage');

        const filter = { key: 'foo', value: 'bar' };
        service.remove(filter);

        expect(websocketService.sendMessage).toHaveBeenCalledWith(
          'removeFilter',
          {
            filter,
          }
        );
      }
    ));
  });

  describe('clear', () => {
    it('should tell backend service to clear filters', inject(
      [LabelFilterService, WebsocketService],
      (service: LabelFilterService, websocketService) => {
        spyOn(websocketService, 'sendMessage');

        service.clear();

        expect(websocketService.sendMessage).toHaveBeenCalledWith(
          'clearFilters',
          {}
        );
      }
    ));
  });

  describe('decodeFilter', () => {
    describe('with valid input', () => {
      it('should return a filter', inject(
        [LabelFilterService],
        (svc: LabelFilterService) => {
          const filter = { key: 'foo', value: 'bar' };
          expect(svc.decodeFilter('foo:bar')).toEqual(filter);
        }
      ));
    });

    describe('with invalid input', () => {
      it('should null', inject(
        [LabelFilterService],
        (svc: LabelFilterService) => {
          expect(svc.decodeFilter('')).toBeNull();
        }
      ));
    });
  });
});
