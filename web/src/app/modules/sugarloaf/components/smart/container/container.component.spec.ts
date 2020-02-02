/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { HttpClientTestingModule } from '@angular/common/http/testing';
import { async, inject, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { ContainerComponent } from './container.component';
import { NamespaceComponent } from '../namespace/namespace.component';
import { PageNotFoundComponent } from '../page-not-found/page-not-found.component';
import { InputFilterComponent } from '../input-filter/input-filter.component';
import { NotifierComponent } from '../notifier/notifier.component';
import { NavigationComponent } from '../navigation/navigation.component';
import { ContextSelectorComponent } from '../../../../shared/components/smart/context-selector/context-selector.component';
import { DefaultPipe } from '../../../../shared/pipes/default/default.pipe';
import { NgSelectModule } from '@ng-select/ng-select';
import {
  BackendService,
  WebsocketService,
} from '../../../../shared/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../../shared/services/websocket/mock';
import { ClarityIcons } from '@clr/icons';
import { ThemeSwitchButtonComponent } from '../theme-switch/theme-switch-button.component';
import { QuickSwitcherComponent } from '../quick-switcher/quick-switcher.component';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: WebsocketService, useClass: WebsocketServiceMock },
        { provide: window, useValue: ClarityIcons },
      ],
      imports: [
        RouterTestingModule,
        ClarityModule,
        HttpClientTestingModule,
        FormsModule,
        NgSelectModule,
      ],
      declarations: [
        ContainerComponent,
        NamespaceComponent,
        PageNotFoundComponent,
        InputFilterComponent,
        NotifierComponent,
        NavigationComponent,
        ContextSelectorComponent,
        DefaultPipe,
        ThemeSwitchButtonComponent,
        QuickSwitcherComponent,
      ],
    }).compileComponents();
  }));

  it('should create the home', () => {
    const fixture = TestBed.createComponent(ContainerComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  describe('at startup', () => {
    it('opens a websocket connection', inject(
      [WebsocketService],
      (websocketService: WebsocketServiceMock) => {
        const fixture = TestBed.createComponent(ContainerComponent);
        fixture.detectChanges();

        expect(websocketService.isOpen).toBeTruthy();
      }
    ));
  });
});
