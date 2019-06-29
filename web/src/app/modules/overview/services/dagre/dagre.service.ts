// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ElementRef, Injectable } from '@angular/core';
import * as d3 from 'd3';
import * as dagreD3 from 'dagre-d3';

@Injectable({
  providedIn: 'root',
})
export class DagreService {
  constructor() {}

  render(parentElement: ElementRef, g) {
    const viewer = parentElement.nativeElement;

    d3.select(viewer)
      .selectAll('*')
      .remove();
    const svg = d3
      .select(viewer)
      .append('svg')
      .attr('width', viewer.offsetWidth)
      .attr('height', viewer.offsetHeight)
      .attr('class', 'dagre-d3');

    const inner = svg.append('g');

    const render = new dagreD3.render();

    render.shapes().record = render.shapes().rect;

    render(inner, g);

    const initialScale = 1.0;

    const svgWidth = parseInt(svg.attr('width'), 10);

    // Set up zoom support
    const zoom = d3.zoom().on('zoom', () => {
      inner.attr('transform', d3.event.transform);
    });
    svg.call(zoom);

    const translateX = (svgWidth - g.graph().width * initialScale) / 2;

    // Center the graph
    const translation = d3.zoomIdentity
      .translate(translateX, 80)
      .scale(initialScale);

    if (!Number.isNaN(translateX)) {
      svg.call(zoom.transform, translation);
    }
  }
}
