// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { inject, TestBed } from '@angular/core/testing';
import { ThemeService } from './theme.service';
import { DOCUMENT } from '@angular/common';
import {
  OverlayScrollbarsComponent,
  OverlayscrollbarsModule,
} from 'overlayscrollbars-ngx';

describe('ThemeService', () => {
  let service: ThemeService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [OverlayScrollbarsComponent],
      providers: [ThemeService, OverlayscrollbarsModule, Document],
    });

    service = TestBed.inject(ThemeService);
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
