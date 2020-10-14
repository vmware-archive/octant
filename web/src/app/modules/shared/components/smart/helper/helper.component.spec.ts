// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { HelperComponent } from './helper.component';
import { HelperService } from '../../../services/helper/helper.service';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { windowProvider, WindowToken } from '../../../../../window';

describe('HelperComponent', () => {
  let component: HelperComponent;
  let fixture: ComponentFixture<HelperComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [HelperComponent],
        providers: [
          { provide: HelperService },
          { provide: WindowToken, useFactory: windowProvider },
        ],
        imports: [BrowserAnimationsModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(HelperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
