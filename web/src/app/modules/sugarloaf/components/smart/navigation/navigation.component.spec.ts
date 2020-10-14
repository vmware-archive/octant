// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { DefaultPipe } from '../../../../shared/pipes/default/default.pipe';
import { NavigationComponent } from './navigation.component';
import { NamespaceComponent } from '../namespace/namespace.component';
import { NgSelectModule } from '@ng-select/ng-select';
import { ThemeSwitchButtonComponent } from '../theme-switch/theme-switch-button.component';
import { MonacoEditorConfig, MonacoProviderService } from 'ng-monaco-editor';
import { windowProvider, WindowToken } from '../../../../../window';

describe('NavigationComponent', () => {
  let component: NavigationComponent;
  let fixture: ComponentFixture<NavigationComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        providers: [
          MonacoProviderService,
          MonacoEditorConfig,
          { provide: WindowToken, useFactory: windowProvider },
        ],
        imports: [NgSelectModule],
        declarations: [
          NavigationComponent,
          NamespaceComponent,
          DefaultPipe,
          ThemeSwitchButtonComponent,
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(NavigationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
