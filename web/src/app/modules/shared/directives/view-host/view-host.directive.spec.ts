import { TestBed } from '@angular/core/testing';
import { EditorComponent } from '../../components/smart/editor/editor.component';
/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { ViewHostDirective } from './view-host.directive';
import { SharedModule } from '../../shared.module';

describe('ViewHostDirective', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [EditorComponent],
      imports: [SharedModule],
    });
  });
  it('should create an instance', () => {
    const directive = new ViewHostDirective(undefined);
    expect(directive).toBeTruthy();
  });
});
