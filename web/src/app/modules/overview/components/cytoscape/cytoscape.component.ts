// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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

import cytoscape, { SingularData, Stylesheet } from 'cytoscape';
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
  @ViewChild('cy') private cy: ElementRef;
  @Input() public elements: any;
  @Input() public style: Stylesheet[];
  @Input() public layout: any;
  @Input() public zoom: any;

  @Output() select: EventEmitter<any> = new EventEmitter<any>();

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
    const cy = cytoscape({
      container: cyContainer,
      layout: this.layout,
      minZoom: this.zoom.min,
      maxZoom: this.zoom.max,
      style: this.style,
      elements: this.elements,
    });

    cy.on('tap', 'node', e => {
      const node: SingularData = e.target;
      localSelect.emit(node.data());
    });
  }
}
