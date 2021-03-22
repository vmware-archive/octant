// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { CytoscapeComponent } from './cytoscape.component';

describe('CytoscapeComponent', () => {
  let component: CytoscapeComponent;
  let fixture: ComponentFixture<CytoscapeComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [CytoscapeComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(CytoscapeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
    expect(component.nodes.length).toEqual(0);
  });

  it('should select first node', done => {
    component.elements = {
      nodes: [
        {
          data: {
            id: '16428c94-a848-47d5-b1e3-c8245b57959b',
            label1: 'metadata-proxy-v0.1',
            label2: 'apps/v1 DaemonSet',
            weight: 100,
            status: 'ok',
            colorCode: '#60b515',
          },
        },
      ],
      edges: [],
    };

    component.render();
    fixture.detectChanges();

    setTimeout(() => {
      expect(component.nodes().length).toEqual(1);
      const node = component.nodes()[0];

      expect(node).not.toBeNull();
      expect(node.id()).toEqual('16428c94-a848-47d5-b1e3-c8245b57959b');
      expect(node.isNode()).toBeTrue();
      expect(node.selected()).toBeTrue();
      done();
    }, 100); // wait for cytoscape to update the view
  });
});
