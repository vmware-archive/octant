// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component } from '@angular/core';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TextView } from '../../../../models/content';
import { TextComponent } from './text.component';

@Component({
  template: '<app-view-text [view]="view"></app-view-text>',
})
class TestWrapperComponent {
  view: TextView;
}

describe('TextComponent', () => {
  describe('handle changes', () => {
    let component: TestWrapperComponent;
    let fixture: ComponentFixture<TestWrapperComponent>;

    beforeEach(async(() => {
      TestBed.configureTestingModule({
        declarations: [TestWrapperComponent, TextComponent],
      }).compileComponents();
    }));

    beforeEach(() => {
      fixture = TestBed.createComponent(TestWrapperComponent);
      component = fixture.componentInstance;
    });

    it('should show text', () => {
      const element: HTMLDivElement = fixture.nativeElement;
      component.view = {
        config: { value: '*text*' },
        metadata: { type: 'text', title: [], accessor: 'accessor' },
      };
      fixture.detectChanges();

      expect(element.querySelector('app-view-text div')).toBeNull();
      expect(element.querySelector('app-view-text').innerHTML).toContain(
        '*text*'
      );
    });

    it('should show markdown text', () => {
      const element: HTMLDivElement = fixture.nativeElement;
      component.view = {
        config: { value: '*text*', isMarkdown: true },
        metadata: { type: 'text', title: [], accessor: 'accessor' },
      };
      fixture.detectChanges();

      expect(element.querySelector('app-view-text markdown')).toBeDefined();
      expect(element.querySelector('app-view-text').innerHTML).toContain(
        '*text*'
      );
    });

    it('should render new markdown text after value has changed', () => {
      component.view = {
        config: { value: '*text*', isMarkdown: true },
        metadata: { type: 'text', title: [], accessor: 'accessor' },
      };
      fixture.detectChanges();

      let element: HTMLElement = fixture.nativeElement;
      expect(element.querySelector('app-view-text div').innerHTML).toEqual(
        '<p><em>text</em></p>\n'
      );

      component.view = {
        config: { value: '# header', isMarkdown: true },
        metadata: { type: 'text', title: [], accessor: 'accessor' },
      };
      fixture.detectChanges();

      element = fixture.nativeElement;
      expect(element.querySelector('app-view-text div').innerHTML).toEqual(
        '<h1 id="header">header</h1>\n'
      );
    });
  });
});
