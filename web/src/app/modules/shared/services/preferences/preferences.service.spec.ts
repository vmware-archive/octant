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
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('collapsed set properly', () => {
    const defaultPrefs = service.getPreferences();
    const defaultElement = defaultPrefs.panels[0].sections[1].elements[0];
    expect(defaultElement.name).toEqual('general.navigation');
    expect(defaultElement.value).toEqual('expanded');

    service.navCollapsed.next(true);
    const newPrefs = service.getPreferences();
    const newElement = newPrefs.panels[0].sections[1].elements[0];
    expect(newElement.name).toEqual('general.navigation');
    expect(newElement.value).toEqual('collapsed');
  });
});
