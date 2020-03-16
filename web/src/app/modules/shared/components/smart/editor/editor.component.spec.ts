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
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
