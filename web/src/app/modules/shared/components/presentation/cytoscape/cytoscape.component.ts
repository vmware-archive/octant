// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  Output,
  Renderer2,
  SimpleChanges,
  ViewChild,
} from '@angular/core';

import cytoscape, { NodeCollection, SingularData, Stylesheet } from 'cytoscape';
import dagre from 'cytoscape-dagre';
import nodeHtmlLabel from 'cytoscape-node-html-label';

cytoscape.use(dagre);
nodeHtmlLabel(cytoscape);

@Component({
  selector: 'app-cytoscape',
  template: '<div #cy class="cy"></div>',
  styleUrls: ['./cytoscape.component.scss'],
})
export class CytoscapeComponent implements OnChanges, OnDestroy {
  @ViewChild('cy', { static: true }) private cy: ElementRef;
  @Input() public elements: any;
  @Input() public style: Stylesheet[];
  @Input() public layout: any;
  @Input() public zoom: any;

  @Output() select: EventEmitter<any> = new EventEmitter<any>();
  @Output() doubleClick: EventEmitter<any> = new EventEmitter<any>();

  private instance: cytoscape.Core;
  private doubleClickDelay = 400;
  private previousTapStamp;

  constructor(private renderer: Renderer2) {
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

  ngOnDestroy(): void {
    if (this.instance && !this.instance.destroyed()) {
      this.instance.destroy();
    }
  }

  public render() {
    const cyContainer = this.renderer.selectRootElement(this.cy.nativeElement);
    const localSelect = this.select;
    const localDoubleClick = this.doubleClick;

    this.instance = cytoscape({
      container: cyContainer,
      layout: this.layout,
      minZoom: this.zoom.min,
      maxZoom: this.zoom.max,
      style: this.style,
      elements: this.elements,
    });

    this.instance.on('tap', 'node', e => {
      const currentTapStamp = e.timeStamp;
      const msFromLastTap = currentTapStamp - this.previousTapStamp;
      const node: SingularData = e.target;

      if (msFromLastTap < this.doubleClickDelay) {
        localDoubleClick.emit(node.data());
      } else {
        localSelect.emit(node.data());
      }
      this.previousTapStamp = currentTapStamp;
    });

    this.instance.one('render', _ => {
      const firstNode = this.instance.nodes().first();
      this.instance.nodes().unselect();
      firstNode.select();
    });

    // @ts-ignore
    this.instance.nodeHtmlLabel([
      {
        query: 'node',
        valign: 'top',
        halign: 'left',
        valignBox: 'bottom',
        halignBox: 'right',
        tpl: data =>
          '<div class="label-header"><p class="label1">' +
          data.label1 +
          '</p>' +
          '<p class="label2">' +
          data.label2 +
          '</p></div>',
      },
    ]);
  }

  public nodes(): NodeCollection {
    return this.instance.nodes();
  }
}
