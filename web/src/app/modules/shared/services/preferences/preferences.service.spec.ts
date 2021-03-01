/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { TestBed } from '@angular/core/testing';

import { PreferencesService } from './preferences.service';

describe('PreferencesService', () => {
  let service: PreferencesService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(PreferencesService);
    service.reset();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('collapsed set properly', () => {
    service.preferences.get('navigation.collapsed').subject.next(false);
    const defaultPrefs = service.getPreferences();
    const defaultElement = defaultPrefs.panels[1].sections[0].elements[0];
    expect(defaultElement.name).toEqual('navigation.collapsed');
    expect(defaultElement.value).toEqual('expanded');

    service.preferences.get('navigation.collapsed').subject.next(true);
    const newPrefs = service.getPreferences();
    const newElement = newPrefs.panels[1].sections[0].elements[0];
    expect(newElement.name).toEqual('navigation.collapsed');
    expect(newElement.value).toEqual('collapsed');
  });

  it('table pagination set properly', () => {
    service.preferences.get('general.pageSize').subject.next(50);
    const defaultPrefs = service.getPreferences();
    const defaultElement = defaultPrefs.panels[0].sections[1].elements[0];
    expect(defaultElement.name).toEqual('general.pageSize');
    expect(defaultElement.value.toString()).toEqual('50');

    service.preferences.get('general.pageSize').subject.next(100);
    const newPrefs = service.getPreferences();
    const newElement = newPrefs.panels[0].sections[1].elements[0];
    expect(newElement.name).toEqual('general.pageSize');
    expect(newElement.value.toString()).toEqual('100');
  });

  it('properties are persisted', () => {
    // default values
    service.preferences.get('navigation.collapsed').subject.next(true);
    service.preferences.get('navigation.labels').subject.next(true);
    service.preferences.get('general.pageSize').subject.next(10);
    expect(
      JSON.parse(service.getStoredValue('navigation.labels', true))
    ).toEqual(true);
    expect(
      JSON.parse(service.getStoredValue('navigation.collapsed', true))
    ).toEqual(true);
    expect(JSON.parse(service.getStoredValue('general.pageSize', 10))).toEqual(
      10
    );

    // change through exposed variables
    service.preferences.get('navigation.collapsed').subject.next(false);
    expect(
      JSON.parse(service.getStoredValue('navigation.collapsed', true))
    ).toEqual(false);

    service.preferences.get('navigation.labels').subject.next(false);
    expect(
      JSON.parse(service.getStoredValue('navigation.labels', true))
    ).toEqual(false);

    service.preferences.get('general.pageSize').subject.next(20);
    expect(JSON.parse(service.getStoredValue('general.pageSize', 10))).toEqual(
      20
    );

    // change by modifying stored values
    service.setStoredValue('navigation.collapsed', true);
    expect(
      JSON.parse(service.getStoredValue('navigation.collapsed', false))
    ).toEqual(true);

    service.setStoredValue('navigation.labels', true);
    expect(
      JSON.parse(service.getStoredValue('navigation.labels', false))
    ).toEqual(true);

    service.setStoredValue('general.pageSize', 100);
    expect(JSON.parse(service.getStoredValue('general.pageSize', 10))).toEqual(
      100
    );
  });
});
