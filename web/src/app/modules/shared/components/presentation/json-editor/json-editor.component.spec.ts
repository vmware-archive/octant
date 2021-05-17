/* Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSONEditorComponent } from './json-editor.component';
import { ComponentFixture, TestBed } from '@angular/core/testing';

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
});
