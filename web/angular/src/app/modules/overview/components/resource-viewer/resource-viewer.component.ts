import { AfterViewChecked, Component, ElementRef, Input, ViewChild, ViewEncapsulation } from '@angular/core';
import * as d3 from 'd3';
import * as dagreD3 from 'dagre-d3';
import { ResourceViewerView } from 'src/app/models/content';
import { Edge } from 'dagre';
import { DagreService } from '../../services/dagre/dagre.service';

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

  constructor(private dagreService: DagreService) {}



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

    this.dagreService.render(this.viewer, g);
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
