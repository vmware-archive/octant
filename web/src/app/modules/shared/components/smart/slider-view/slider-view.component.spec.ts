/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SliderViewComponent } from './slider-view.component';
import { SharedModule } from '../../../shared.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { windowProvider, WindowToken } from '../../../../../window';

describe('SliderViewComponent', () => {
  let component: SliderViewComponent;
  let fixture: ComponentFixture<SliderViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [SharedModule, BrowserAnimationsModule],
      providers: [{ provide: WindowToken, useFactory: windowProvider }],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SliderViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('slide', () => {
    it('triggers animationState when called', () => {
      expect(component.animationState).toEqual('out');
      component.slide();
      expect(component.animationState).toEqual('in');
      component.slide();
      expect(component.animationState).toEqual('out');
    });
  });
});
