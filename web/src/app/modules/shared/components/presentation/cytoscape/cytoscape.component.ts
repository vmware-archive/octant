// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  Renderer2,
  SimpleChanges,
  ViewChild,
} from '@angular/core';

import cytoscape, { NodeSingular, SingularData, Stylesheet } from 'cytoscape';
import { hideChildren, positionChildren } from './octant.layout';
import dagre from 'cytoscape-dagre';

cytoscape.use(dagre);

@Component({
  selector: 'app-cytoscape',
  template: '<div #cy class="cy"></div>',
  styles: [
    `
      .cy {
        height: 100%;
        width: 100%;
        position: relative;
        left: 0;
        top: 0;
      }
    `,
  ],
})
export class CytoscapeComponent implements OnChanges {
  @ViewChild('cy', { static: true }) private cy: ElementRef;
  @Input() public elements: any;
  @Input() public style: Stylesheet[];
  @Input() public layout: any;
  @Input() public zoom: any;

  @Output() select: EventEmitter<any> = new EventEmitter<any>();

  cytoscape: cytoscape.Core;
  applied = false;
  moveStarted = false;

  constructor(private renderer: Renderer2, private el: ElementRef) {
    this.layout = this.layout || {
      name: 'grid',
      directed: true,
    };

    this.zoom = this.zoom || {
      min: 0.1,
      max: 1.5,
    };
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.render();
  }

  public render() {
    const cyContainer = this.renderer.selectRootElement(this.cy.nativeElement);
    const localSelect = this.select;
    const options: cytoscape.CytoscapeOptions = {
      container: cyContainer,
      layout: this.layout,
      minZoom: this.zoom.min,
      maxZoom: this.zoom.max,
      style: this.style,
      elements: this.elements,
    };
    this.cytoscape = cytoscape(options);

    this.cytoscape.on('tap', 'node', e => {
      const node: SingularData = e.target;
      localSelect.emit(node.data());
    });

    this.cytoscape.on('layoutstop', e => {
      if (!this.applied) {
        this.applied = true;
        this.cytoscape
          .nodes('[hasChildren]')
          .forEach(node => positionChildren(this.cytoscape, node));
        this.cytoscape.fit();
      }
    });

    this.cytoscape.on('drag', 'node', e => {
      const node: NodeSingular = e.target;
      if (!this.moveStarted) {
        this.moveStarted = true;
        hideChildren(this.cytoscape, node);
      }
    });

    this.cytoscape.on('dragfree', 'node', e => {
      const node: NodeSingular = e.target;
      positionChildren(this.cytoscape, node);
      this.moveStarted = false;
    });
  }
}
