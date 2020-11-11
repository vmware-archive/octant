// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { async, TestBed } from '@angular/core/testing';
import { EditorComponent } from '../../components/smart/editor/editor.component';
import { HighlightModule, HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';
import { IconService } from './icon.service';
import { SharedModule } from '../../shared.module';

describe('IconService', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [HighlightModule, SharedModule],
      declarations: [EditorComponent],
      providers: [
        {
          provide: HIGHLIGHT_OPTIONS,
          useValue: {
            languages: {
              json: () => import('highlight.js/lib/languages/json'),
              yaml: () => import('highlight.js/lib/languages/yaml'),
            },
          },
        },
      ],
    });
  }));

  it('should be created', () => {
    const service: IconService = TestBed.inject(IconService);
    expect(service).toBeTruthy();
  });
});
