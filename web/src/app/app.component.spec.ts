// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { HttpClientTestingModule } from '@angular/common/http/testing';
import { async, inject, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { AppComponent } from './app.component';
import { NamespaceComponent } from './components/namespace/namespace.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { InputFilterComponent } from './components/input-filter/input-filter.component';
import { NotifierComponent } from './components/notifier/notifier.component';
import { NavigationComponent } from './components/navigation/navigation.component';
import { ContextSelectorComponent } from './modules/overview/components/context-selector/context-selector.component';
import { DefaultPipe } from './modules/overview/pipes/default.pipe';
import { NgSelectModule } from '@ng-select/ng-select';
import {
  BackendService,
  WebsocketService,
} from './modules/overview/services/websocket/websocket.service';
import { WebsocketServiceMock } from './modules/overview/services/websocket/mock';
import { ClarityIcons } from '@clr/icons';
import { ThemeSwitchButtonComponent } from './modules/overview/components/theme-switch/theme-switch-button.component';
import { QuickSwitcherComponent } from './components/quick-switcher/quick-switcher.component';

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
        AppComponent,
        NamespaceComponent,
        PageNotFoundComponent,
        InputFilterComponent,
        NotifierComponent,
        NavigationComponent,
        ContextSelectorComponent,
        DefaultPipe,
        ThemeSwitchButtonComponent,
        QuickSwitcherComponent
      ],
    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  describe('at startup', () => {
    it('opens a websocket connection', inject(
      [WebsocketService],
      (websocketService: WebsocketServiceMock) => {
        const fixture = TestBed.createComponent(AppComponent);
        fixture.detectChanges();

        expect(websocketService.isOpen).toBeTruthy();
      }
    ));
  });
});
