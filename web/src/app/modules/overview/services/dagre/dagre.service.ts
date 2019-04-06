import { Injectable, ElementRef } from '@angular/core';
import * as d3 from 'd3';
import * as dagreD3 from 'dagre-d3';
import { Graph } from 'graphlib';

@Injectable({
  providedIn: 'root'
})
export class DagreService {
  constructor() { }

  render(parentElement: ElementRef, g) {
    const viewer = parentElement.nativeElement;

    d3.select(viewer).selectAll('*').remove();
    const svg = d3.select(viewer).append('svg')
      .attr('width', viewer.offsetWidth)
      .attr('height', viewer.offsetHeight)
      .attr('class', 'dagre-d3');

    const inner = svg.append('g');

    const render = new dagreD3.render();

    render.shapes().record = render.shapes().rect;

    render(inner, g);

    const initialScale = 1.2;

    const width = parseInt(svg.attr('width'), 10);
    const height = parseInt(svg.attr('height'), 10);

    // Set up zoom support
    const zoom = d3.zoom()
      .on('zoom', () => {
        inner.attr('transform', d3.event.transform);
      });
    svg.call(zoom);

    // Center the graph
    const translation = d3.zoomIdentity.translate(
      (width - g.graph().width * initialScale) / 2,
      (height - g.graph().height * initialScale) / 2,
      ).scale(initialScale);
    svg.call(zoom.transform, translation);
  }
}
