/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { HttpClientTestingModule } from '@angular/common/http/testing';
import { async, inject, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { ClarityModule, ClrPopoverToggleService } from '@clr/angular';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ContainerComponent } from './container.component';
import { NamespaceComponent } from '../namespace/namespace.component';
import { PageNotFoundComponent } from '../page-not-found/page-not-found.component';
import { PreferencesComponent } from '../../../../shared/components/presentation/preferences/preferences.component';
import { HelperComponent } from '../../../../shared/components/smart/helper/helper.component';
import { InputFilterComponent } from '../input-filter/input-filter.component';
import { NotifierComponent } from '../notifier/notifier.component';
import { NavigationComponent } from '../navigation/navigation.component';
import { ContextSelectorComponent } from '../../../../shared/components/smart/context-selector/context-selector.component';
import { DefaultPipe } from '../../../../shared/pipes/default/default.pipe';
import { FilterTextPipe } from '../../../pipes/filtertext/filtertext.pipe';
import { NgSelectModule } from '@ng-select/ng-select';
import { WebsocketService } from '../../../../shared/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../../shared/services/websocket/mock';
import { ClarityIcons } from '@clr/icons';
import { ThemeSwitchButtonComponent } from '../theme-switch/theme-switch-button.component';
import { QuickSwitcherComponent } from '../quick-switcher/quick-switcher.component';
import { MonacoEditorConfig, MonacoProviderService } from 'ng-monaco-editor';
import { UploaderComponent } from '../uploader/uploader.component';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: WebsocketService, useClass: WebsocketServiceMock },
        { provide: window, useValue: ClarityIcons },
        ClrPopoverToggleService,
        MonacoProviderService,
        MonacoEditorConfig,
      ],
      imports: [
        BrowserModule,
        RouterTestingModule,
        ClarityModule,
        HttpClientTestingModule,
        FormsModule,
        NgSelectModule,
        ReactiveFormsModule,
      ],
      declarations: [
        ContainerComponent,
        NamespaceComponent,
        PageNotFoundComponent,
        HelperComponent,
        PreferencesComponent,
        InputFilterComponent,
        NotifierComponent,
        NavigationComponent,
        ContextSelectorComponent,
        DefaultPipe,
        FilterTextPipe,
        ThemeSwitchButtonComponent,
        QuickSwitcherComponent,
        UploaderComponent,
      ],
    }).compileComponents();
  }));

  it('should create the home', () => {
    const fixture = TestBed.createComponent(ContainerComponent);
    const app = fixture.debugElement.componentInstance;
    fixture.detectChanges();
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
