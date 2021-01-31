/* Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { EditorComponent } from './editor.component';
import { MonacoEditorModule } from '@materia-ui/ngx-monaco-editor';
import { windowProvider, WindowToken } from '../../../../../window';

describe('EditorComponent', () => {
  let component: EditorComponent;
  let fixture: ComponentFixture<EditorComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
        imports: [MonacoEditorModule],
        declarations: [EditorComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(EditorComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display the editor', () => {
    const root: HTMLElement = fixture.nativeElement;
    const editorElement: SVGPathElement = root.querySelector(
      '.editor-container .editor'
    );
    expect(editorElement).not.toBeNull();
    expect(editorElement.classList.contains('editor')).toBeTruthy();
  });
});
