import { AfterViewChecked, Component, ElementRef, Input, ViewChild, ViewEncapsulation } from '@angular/core';
import * as d3 from 'd3';
import * as dagreD3 from 'dagre-d3';
import { ResourceViewerView } from 'src/app/models/content';
import { Edge } from 'dagre';

interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
}

class ResourceNode {
  constructor(
    private readonly id: string,
    private readonly object: ResourceObject,
    private readonly isSelected: boolean
  ) {}

  toDescriptor(): dagreD3.Label {
    let nodeClass = `node-${this.object.status}`;
    if (this.isSelected) {
      nodeClass += ` selected`;
    }

    return {
      id: this.id,
      label: `${this.title()}${this.subTitle()}`,
      labelType: 'html',
      class: `${nodeClass}`,
    };
  }

  title(): string {
    return `<div class="resource-name">${this.object.name}</div>`;
  }

  subTitle(): string {
    return `<div class="resource-type">${this.object.apiVersion} ${this.object.kind}</div>`;
  }
}



@Component({
  selector: 'app-view-resource-viewer',
  template: `
    <div class="resourceViewer" #viewer>
    </div>
  `,
  styleUrls: ['./resource-viewer.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ResourceViewerComponent implements AfterViewChecked {

  @ViewChild('viewer') private viewer: ElementRef;

  @Input() view: ResourceViewerView;

  constructor() {}



  ngAfterViewChecked(): void {
    this.updateGraph();
  }

  updateGraph() {
    const viewer = this.viewer.nativeElement;
    if (viewer.offsetWidth === 0 || viewer.offsetHeight === 0) {
      // nothing to do until the viewer has dimensions
      return;
    }

    const g = new dagreD3.graphlib.Graph().setGraph({});

    for (const [id, label] of Object.entries(this.nodes())) {
      g.setNode(id, label);
    }

    g.nodes().forEach((v: any) => {
      const node = g.node(v);
      node.rx = node.ry = 4;
    });

    this.edges().forEach((edge) => g.setEdge(edge[0], edge[1], edge[2]));


    d3.select(viewer).selectAll('*').remove();
    const svg = d3.select(viewer).append('svg')
      .attr('width', viewer.offsetWidth)
      .attr('height', viewer.offsetHeight)
      .attr('class', 'dagre-d3');

    const inner = svg.append('g');

    const render = new dagreD3.render();
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

  edges(): Array<Edge> {
    const adjacencyList = this.view.config.edges;
    const edges: Array<Edge> = [];

    if (adjacencyList) {
      for (const [node, nodeEdges] of Object.entries(adjacencyList)) {
        edges.push(
          ...nodeEdges.map((e) => [
            node,
            e.node,
            {
              arrowhead: 'undirected',
            },
          ])
        );
      }
    }

    return edges;
  }

  nodes() {
    const objects = this.view.config.nodes;

    const nodes: { [key: string]: dagreD3.Label } = {};

    for (const [id, object] of Object.entries(objects)) {
      nodes[id] = new ResourceNode(id, object, false).toDescriptor();
    }

    return nodes;
  }
}
