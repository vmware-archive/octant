/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { HttpClientTestingModule } from '@angular/common/http/testing';
import {
  async,
  ComponentFixture,
  inject,
  TestBed,
} from '@angular/core/testing';
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
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../../../data/services/websocket/mock';
import { ClarityIcons } from '@clr/icons';
import { ThemeSwitchButtonComponent } from '../theme-switch/theme-switch-button.component';
import { QuickSwitcherComponent } from '../quick-switcher/quick-switcher.component';
import {
  MonacoEditorConfig,
  MonacoEditorModule,
  MonacoProviderService,
} from 'ng-monaco-editor';
import { UploaderComponent } from '../uploader/uploader.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { windowProvider, WindowToken } from '../../../../../window';
import { SharedModule } from 'src/app/modules/shared/shared.module';
import { ApplyYAMLComponent } from '../apply-yaml/apply-yaml.component';
import { EditorComponent } from 'src/app/modules/shared/components/smart/editor/editor.component';

describe('AppComponent', () => {
  let component: ContainerComponent;
  let fixture: ComponentFixture<ContainerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: WebsocketService, useClass: WebsocketServiceMock },
        { provide: WindowToken, useFactory: windowProvider },
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
        BrowserAnimationsModule,
        SharedModule,
      ],
      declarations: [
        ApplyYAMLComponent,
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

  beforeEach(() => {
    fixture = TestBed.createComponent(ContainerComponent);
    component = fixture.componentInstance;
  });

  afterEach(() => {
    TestBed.resetTestingModule();
  });

  it('should create the home', () => {
    expect(component).toBeTruthy();
  });

  describe('at startup', () => {
    it('opens a websocket connection', inject(
      [WebsocketService],
      (websocketService: WebsocketServiceMock) => {
        fixture.detectChanges();
        expect(websocketService.isOpen).toBeTruthy();
      }
    ));
  });
});
