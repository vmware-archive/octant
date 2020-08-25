/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewContainerComponent } from './view-container.component';
import { DYNAMIC_COMPONENTS_MAPPING } from '../../dynamic-components';
import { TextComponent } from '../presentation/text/text.component';

describe('ViewContainerComponent', () => {
  let component: ViewContainerComponent;
  let fixture: ComponentFixture<ViewContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ViewContainerComponent],
      providers: [
        {
          provide: DYNAMIC_COMPONENTS_MAPPING,
          useValue: {
            text: TextComponent,
          },
        },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ViewContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
