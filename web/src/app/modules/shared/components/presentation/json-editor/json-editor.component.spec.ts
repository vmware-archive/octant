/* Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSONEditorComponent } from './json-editor.component';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { DebugElement } from '@angular/core';
import { By } from '@angular/platform-browser';
import { JSONEditorView } from '../../../models/content';

describe('JSONEditorComponent', () => {
  let component: JSONEditorComponent;
  let fixture: ComponentFixture<JSONEditorComponent>;

  beforeEach(() => {
    fixture = TestBed.createComponent(JSONEditorComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('displays json editor', () => {
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('.jsoneditor')).toBeTruthy();
  });

  it('parses invalid json', () => {
    component.content = '{123';
    const error = { error: 'cannot parse json' };
    const result = component.isValidJson(component.content);
    expect(result).toEqual(error);
  });

  it('parses valid json', () => {
    component.content = { hello: 'world' };
    const result = component.isValidJson(component.content);
    expect(result).toEqual(component.content);
  });

  it('respects collapsed flag', () => {
    component.view = {
      config: {
        mode: 'view',
        content: '{ "hello": "world", "my": "world" }',
        collapsed: true,
      },
      metadata: {
        type: 'jsonEditor',
      },
    } as JSONEditorView;

    fixture.detectChanges();

    const editorDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.jsoneditor')
    );
    const editorNativeElement: HTMLDivElement =
      editorDebugElement.nativeElement;

    expect(editorNativeElement.clientHeight).toBeLessThan(50);
  });
});
