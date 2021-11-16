// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { GraphvizComponent } from './graphviz.component';
import { Component } from '@angular/core';
import { GraphvizView } from '../../../models/content';

@Component({
  template: '<app-view-graphviz [view]="view"></app-view-graphviz>',
})
class TestWrapperComponent {
  view: GraphvizView;
}

describe('GraphvizComponent', () => {
  let component: TestWrapperComponent;
  let fixture: ComponentFixture<TestWrapperComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [TestWrapperComponent, GraphvizComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TestWrapperComponent);
    component = fixture.componentInstance;
    component.view = {
      config: { dot: 'digraph {a -> b}' },
      metadata: { type: 'graphviz' },
    };
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show graph', () => {
    const element: HTMLDivElement = fixture.nativeElement;
    fixture.detectChanges();

    expect(element.querySelector('app-view-graphviz div')).not.toBeNull();
    expect(element.querySelector('app-view-graphviz div svg g')).not.toBeNull();
    expect(
      element.querySelector('app-view-graphviz div svg g').children.length
    ).toEqual(5);
  });
});
