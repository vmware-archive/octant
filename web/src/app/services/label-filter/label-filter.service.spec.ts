// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TestBed, fakeAsync, tick, async } from '@angular/core/testing';
import { LabelFilterService, Filter } from './label-filter.service';
import { Router } from '@angular/router';
import { NgZone } from '@angular/core';

describe('LabelFilterService', () => {
  let service: LabelFilterService;
  let router: Router;
  let ngZone: NgZone;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.get(LabelFilterService);
    router = TestBed.get(Router);
    ngZone = TestBed.get(NgZone);
  });

  it('should be created with no filters', () => {
    expect(service).toBeTruthy();
    const filters = service.filters.getValue();
    expect(filters).toEqual([]);
  });

  it('should add a filter and trigger router', fakeAsync(() => {
    ngZone.run(() => {
      const testFilter: Filter = { key: 'test1', value: 'value1' };
      service.add(testFilter);
      tick();
      const observedFilters = service.filters.getValue();
      expect(observedFilters).toEqual([{ key: 'test1', value: 'value1' }]);
      expect(router.url).toMatch(/\?filter=test1:value1$/i);
    });
  }));

  it('should decode a filter query param', () => {
    const filterQueryParam = 'test1:value1';
    const filter = service.decodeFilter(filterQueryParam);
    expect(filter).toEqual({ key: 'test1', value: 'value1' });
  });

  it('should delete a filter without removing other filters', fakeAsync(() => {
    ngZone.run(() => {
      service.filters.next([
        { key: 'test1', value: 'value1' },
        { key: 'test2', value: 'value2' },
        { key: 'test3', value: 'value3' },
        { key: 'test4', value: 'value4' },
      ]);
      service.remove({ key: 'test3', value: 'value3' });
      tick();
      const observedFilters = service.filters.getValue();
      const expectedFilters: Filter[] = [
        { key: 'test1', value: 'value1' },
        { key: 'test2', value: 'value2' },
        { key: 'test4', value: 'value4' },
      ];
      expect(observedFilters).toEqual(expectedFilters);
      expect(router.url).toMatch(/filter=test1:value1/i);
      expect(router.url).toMatch(/filter=test2:value2/i);
      expect(router.url).toMatch(/filter=test4:value4/i);
      expect(router.url).not.toMatch(/filter=test3:value3/i);
    });
  }));

  it('should clear all filters', fakeAsync(() => {
    ngZone.run(() => {
      service.filters.next([
        { key: 'test1', value: 'value1' },
        { key: 'test2', value: 'value2' },
        { key: 'test3', value: 'value3' },
        { key: 'test4', value: 'value4' },
      ]);
      service.clearAll();
      tick();
      expect(service.filters.getValue()).toEqual([]);
      expect(router.url).toMatch(/\/\?/i);
    });
  }));
});
