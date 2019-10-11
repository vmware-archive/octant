// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { TestBed, inject } from '@angular/core/testing';
import { ThemeService } from './theme-switch.service';
import { DOCUMENT } from '@angular/common';

describe('ThemeService', () => {
  let service: ThemeService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [ThemeService, Document],
    });

    service = TestBed.get(ThemeService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should load light theme file correctly', inject(
    [DOCUMENT],
    (document: Document) => {
      service.loadCSS('assets/css/clr-ui.min.css');

      const themeLink = document.getElementById(
        'client-theme'
      ) as HTMLLinkElement;
      expect(themeLink.href).toContain('assets/css/clr-ui.min.css');
    }
  ));

  it('should load dark theme file correctly', inject(
    [DOCUMENT],
    (document: Document) => {
      service.loadCSS('assets/css/clr-ui-dark.min.css');

      const themeLink = document.getElementById(
        'client-theme'
      ) as HTMLLinkElement;
      expect(themeLink.href).toContain('assets/css/clr-ui-dark.min.css');
    }
  ));
});
