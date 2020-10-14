import { TestBed } from '@angular/core/testing';
import { EditorComponent } from '../../components/smart/editor/editor.component';
// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { DefaultPipe } from './default.pipe';
import { SharedModule } from '../../shared.module';

describe('DefaultPipe', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [EditorComponent],
      providers: [SharedModule],
    });
  });
  it('create an instance', () => {
    const pipe = new DefaultPipe();
    expect(pipe).toBeTruthy();
  });

  it('acts as an identity function with a non empty string', () => {
    const pipe = new DefaultPipe();
    expect(pipe.transform('foo', 'bar')).toEqual('foo');
  });

  it('returns the default with an empty string', () => {
    const pipe = new DefaultPipe();
    expect(pipe.transform('', 'bar')).toEqual('bar');
  });
});
