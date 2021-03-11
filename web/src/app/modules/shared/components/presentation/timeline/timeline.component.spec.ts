// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Component } from '@angular/core';
import { TimelineView } from '../../../models/content';
import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { TimelineComponent } from './timeline.component';
@Component({
  template: '<app-view-timeline [view]="view"></app-view-timeline>',
})
class TestWrapperComponent {
  view: TimelineView;
}

describe('TimelineComponent', () => {
  describe('handle changes', () => {
    let component: TestWrapperComponent;
    let fixture: ComponentFixture<TestWrapperComponent>;

    beforeEach(
      waitForAsync(() => {
        TestBed.configureTestingModule({
          providers: [],
          declarations: [TestWrapperComponent, TimelineComponent],
        }).compileComponents();
      })
    );

    beforeEach(() => {
      fixture = TestBed.createComponent(TestWrapperComponent);
      component = fixture.componentInstance;
    });

    it('should show step', () => {
      const element: HTMLDivElement = fixture.nativeElement;
      component.view = {
        config: {
          steps: [
            {
              state: 'current',
              header: 'header',
              title: 'title',
              description: 'description',
            },
          ],
          vertical: false,
        },
        metadata: { type: 'timeline', title: [], accessor: 'accessor' },
      };
      fixture.detectChanges();

      expect(element.querySelector('app-view-timeline').innerHTML).toContain(
        'description'
      );
    });
  });
});
