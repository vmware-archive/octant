/* Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { EditorComponent } from './editor.component';
import {
  MonacoEditorConfig,
  MonacoEditorModule,
  MonacoProviderService,
} from 'ng-monaco-editor';
import { SharedModule } from '../../../shared.module';

describe('EditorComponent', () => {
  let component: EditorComponent;
  let fixture: ComponentFixture<EditorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      providers: [MonacoProviderService, MonacoEditorConfig],
      imports: [
        MonacoEditorModule.forRoot({
          baseUrl: '',
          defaultOptions: {},
        }),
      ],
      declarations: [EditorComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
